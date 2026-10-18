package hub

import (
	"encoding/json"
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 64 << 10
)

type client struct {
	hub  *Hub
	conn *websocket.Conn
	out  chan []byte
	id   types.UserID
	room *room.Room
	mu   sync.Mutex
}

type Hub struct {
	reg *registry.Registry
	br  *bridge.Bridge
	StatsRecorder room.StatsRecorder

	mu      sync.Mutex
	clients map[*client]struct{}
	rooms   map[types.RoomID]*room.Room

	done chan struct{}
	once sync.Once
}

func New(reg *registry.Registry) *Hub {
	return &Hub{
		reg:     reg,
		clients: make(map[*client]struct{}),
		rooms:   make(map[types.RoomID]*room.Room),
		done:    make(chan struct{}),
	}
}

func (h *Hub) SetBridge(br *bridge.Bridge) { h.br = br }

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	w := &wireRoom{h: h}
	rm, err := h.reg.Open(registry.Spec{
		GameID:    gameID,
		Players:   players,
		Transport: w,
		StatsRecorder: h.StatsRecorder,
	})
	if err != nil {
		return nil, err
	}
	h.mu.Lock()
	h.rooms[rm.ID()] = rm
	h.mu.Unlock()
	return rm, nil
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := types.UserID(r.URL.Query().Get("user"))
	if uid == "" {
		http.Error(w, "user required", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("hub: accept: %v", err)
		return
	}
	conn.SetReadLimit(maxMessageSize)

	c := &client{hub: h, conn: conn, out: make(chan []byte, 64), id: uid}

	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	go c.writeLoop()
	c.readLoop()
}

func (h *Hub) broadcastToRoom(rm *room.Room, frame []byte) {
	h.mu.Lock()
	clients := make([]*client, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.Unlock()
	for _, c := range clients {
		c.mu.Lock()
		my := c.room
		c.mu.Unlock()
		if my != nil && my.ID() == rm.ID() {
			c.enqueue(frame)
		}
	}
	if h.br != nil {
		h.br.PushRaw(frame)
	}
}

func (h *Hub) Shutdown(ctx context.Context) error {
	h.once.Do(func() {
		close(h.done)
		h.mu.Lock()
		for c := range h.clients {
			_ = c.conn.CloseNow()
		}
		h.mu.Unlock()
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func (c *client) readLoop() {
	ctx := context.Background()
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			break
		}
		
		// --- تغییر به پروتکل باینری ---
		m, err := wire.NewInboundMsg(data)
		if err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid binary frame"))
			continue
		}
		// ------------------------------

		c.handle(m)
	}
	c.hub.mu.Lock()
	delete(c.hub.clients, c)
	c.hub.mu.Unlock()
	c.mu.Lock()
	rm := c.room
	c.room = nil
	c.mu.Unlock()
	if rm != nil {
		_ = rm.Leave(ctx, c.id)
	}
}

func (c *client) handle(m wire.InboundMsg) {
	switch m.Type {
	case wire.TypePing:
		c.enqueue(wire.OutboundMsg{Type: wire.TypePong}.Bytes())
	case wire.TypeJoin:
		c.hub.mu.Lock()
		rm, ok := c.hub.rooms[types.RoomID(m.RoomID)]
		c.hub.mu.Unlock()
		if !ok {
			c.enqueue(wire.ErrFrame(wire.CodeNoRoom, "no such room"))
			return
		}
		c.mu.Lock()
		c.room = rm
		c.mu.Unlock()
		{
			jctx, jcancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := rm.Join(jctx, c.id); err != nil {
				jcancel()
				c.enqueue(errFrameFor(err))
				return
			}
			jcancel()
		}
		st, rev := rm.GetState()
		c.enqueue(wire.OutboundMsg{
			Type:   wire.TypeJoined,
			RoomID: string(rm.ID()),
			GameID: string(rm.GameID()),
			Rev:    rev,
			State:  st,
		}.Bytes())
	case wire.TypeLeave:
		c.mu.Lock()
		rm := c.room
		c.room = nil
		c.mu.Unlock()
		if rm != nil {
			_ = rm.Leave(context.Background(), c.id)
		}
	case wire.TypeChat:
		c.mu.Lock(); rm := c.room; c.mu.Unlock()
		if rm == nil { return }
		var p struct{ Text string `json:"text"` }
		if err := json.Unmarshal(m.Payload, &p); err != nil { return }
		if len(p.Text) > 500 { p.Text = p.Text[:500] }
		if p.Text == "" { return }
		c.hub.broadcastToRoom(rm, wire.OutboundMsg{
			Type: wire.TypeChat, RoomID: string(rm.ID()), From: string(c.id),
			Payload: map[string]any{"text": p.Text, "user_id": string(c.id)},
		}.Bytes())
	case wire.TypeEmoji:
		c.mu.Lock(); rm := c.room; c.mu.Unlock()
		if rm == nil { return }
		var p struct{ Emoji string `json:"emoji"` }
		if err := json.Unmarshal(m.Payload, &p); err != nil { return }
		if len(p.Emoji) > 8 { p.Emoji = p.Emoji[:8] }
		if p.Emoji == "" { return }
		c.hub.broadcastToRoom(rm, wire.OutboundMsg{
			Type: wire.TypeEmoji, RoomID: string(rm.ID()), From: string(c.id),
			Payload: map[string]any{"emoji": p.Emoji, "user_id": string(c.id)},
		}.Bytes())
	case wire.TypeState:
		c.mu.Lock()
		rm := c.room
		c.mu.Unlock()
		if rm == nil {
			c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "not joined"))
			return
		}
		st, rev := rm.GetState()
		c.enqueue(wire.OutboundMsg{
			Type:   wire.TypeState,
			RoomID: string(rm.ID()),
			GameID: string(rm.GameID()),
			Rev:    rev,
			State:  st,
		}.Bytes())
	default:
		c.mu.Lock()
		rm := c.room
		c.mu.Unlock()
		if rm == nil {
			c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "join first"))
			return
		}
		cmd, err := m.Command()
		if err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, err.Error()))
			return
		}
		if cmd.RoomID == "" {
			cmd.RoomID = rm.ID()
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := rm.Do(ctx, room.Request{From: c.id, Cmd: cmd}); err != nil {
			c.enqueue(errFrameFor(err))
		}
	}
}

func errFrameFor(err error) []byte {
	switch {
	case errors.Is(err, room.ErrStaleVersion):
		return wire.ErrFrame(wire.CodeStaleVersion, err.Error())
	case errors.Is(err, room.ErrWrongPhase):
		return wire.ErrFrame(wire.CodeWrongPhase, err.Error())
	case errors.Is(err, room.ErrBusy):
		return wire.ErrFrame(wire.CodeBusy, err.Error())
	case errors.Is(err, room.ErrClosed):
		return wire.ErrFrame(wire.CodeServerClosed, err.Error())
	default:
		return wire.ErrFrame(wire.CodeRejected, err.Error())
	}
}

func (c *client) enqueue(frame []byte) {
	select {
	case c.out <- frame:
	default:
	}
}

func (c *client) writeLoop() {
	defer c.conn.CloseNow()
	ctx := context.Background()
	for {
		select {
		case <-c.hub.done:
			return
		case frame := <-c.out:
			wctx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(wctx, websocket.MessageBinary, frame) // ارسال باینری
			cancel()
			if err != nil {
				return
			}
		}
	}
}

type wireRoom struct{ h *Hub }

func (w *wireRoom) Broadcast(frame []byte) {
	if w.h.br != nil {
		w.h.br.PushRaw(frame)
	}
	w.h.mu.Lock()
	clients := make([]*client, 0, len(w.h.clients))
	for c := range w.h.clients {
		clients = append(clients, c)
	}
	w.h.mu.Unlock()
	for _, c := range clients {
		c.enqueue(frame)
	}
}

func (w *wireRoom) SendTo(id types.UserID, frame []byte) {
	w.h.mu.Lock()
	clients := make([]*client, 0, 1)
	for c := range w.h.clients {
		if c.id == id {
			clients = append(clients, c)
		}
	}
	w.h.mu.Unlock()
	for _, c := range clients {
		c.enqueue(frame)
	}
}

func (w *wireRoom) CloseSlowClient(id types.UserID, code string, reason string) {
	w.h.mu.Lock()
	var target *client
	for c := range w.h.clients {
		if c.id == id {
			target = c
			break
		}
	}
	w.h.mu.Unlock()
	if target == nil {
		return
	}
	target.enqueue(wire.ErrFrame(code, reason))
	_ = target.conn.CloseNow()
}

func (w *wireRoom) SignalBackpressure(code string, reason string) {
	w.h.mu.Lock()
	clients := make([]*client, 0, len(w.h.clients))
	for c := range w.h.clients {
		clients = append(clients, c)
	}
	w.h.mu.Unlock()
	frame := wire.ErrFrame(code, reason)
	for _, c := range clients {
		c.enqueue(frame)
	}
}
