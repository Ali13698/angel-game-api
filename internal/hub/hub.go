package hub

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

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

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	w := &wireRoom{h: h}
	rm, err := h.reg.Open(registry.Spec{
		GameID:    gameID,
		Players:   players,
		Transport: w,
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
		var m wire.InboundMsg
		if err := json.Unmarshal(data, &m); err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid json"))
			continue
		}
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
			err := c.conn.Write(wctx, websocket.MessageText, frame)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

type wireRoom struct{ h *Hub }

func (w *wireRoom) Broadcast(frame []byte) {
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

// CloseSlowClient در سرریز fanout برای ToPlayer فراخوانی می‌شود.
// یک فریم خطا enqueue می‌کند و سپس کانال را با StatusPolicyViolation می‌بندد.
// غیرمسدودکننده: enqueue و CloseNow آنی برمی‌گردند.
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

// SignalBackpressure در سرریز fanout برای ToRoom فراخوانی می‌شود.
// یک فریم سیگنال به همهٔ کلاینت‌های اتاق enqueue می‌کند (غیرمسدودکننده).
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
