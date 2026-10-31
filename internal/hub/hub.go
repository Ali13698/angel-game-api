package hub

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/room"
	"github.com/coder/websocket"
)

type Hub struct {
	reg            *registry.Registry
	br             *bridge.Bridge
	StatsRecorder  room.StatsRecorder
	EventPublisher room.EventPublisher
	mu             sync.Mutex
	clients        map[*client]struct{}
	rooms          map[types.RoomID]*room.Room
	allowed        map[*client]types.RoomID
	play           *playLobby
	done           chan struct{}
	once           sync.Once
}

func New(reg *registry.Registry, pub room.EventPublisher) *Hub {
	return &Hub{
		reg:            reg,
		EventPublisher: pub,
		clients:        make(map[*client]struct{}),
		rooms:          make(map[types.RoomID]*room.Room),
		allowed:        make(map[*client]types.RoomID),
		play:           newPlayLobby(),
		done:           make(chan struct{}),
	}
}

func (h *Hub) SetBridge(br *bridge.Bridge) {
	h.br = br
}

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	rm, err := h.reg.Open(registry.Spec{
		GameID:         gameID,
		Players:        players,
		Transport:      &wireRoom{h: h},
		StatsRecorder:  h.StatsRecorder,
		EventPublisher: h.EventPublisher,
	})
	if err != nil {
		return nil, err
	}

	h.mu.Lock()
	h.rooms[rm.ID()] = rm
	h.mu.Unlock()

	logger.Info("Room opened", "room_id", rm.ID(), "game_id", gameID)

	for _, player := range players {
		if strings.HasPrefix(string(player), "bot") {
			go h.runBot(rm, types.UserID("bot"), gameID, 3*time.Second)
			break
		}
	}

	return rm, nil
}

func browserOriginOK(r *http.Request) bool {
	u, err := url.Parse(r.Header.Get("Origin"))
	return err == nil &&
		u.Host == r.Host &&
		u.User == nil &&
		(u.Scheme == "https" || u.Scheme == "http")
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !browserOriginOK(r) {
		http.Error(w, "origin rejected", http.StatusForbidden)
		return
	}

	if r.URL.Query().Get("lobby") == "1" {
		h.servePlayLobby(w, r)
		return
	}

	var ticket string
	for _, protocol := range strings.Split(strings.Join(r.Header.Values("Sec-WebSocket-Protocol"), ","), ",") {
		protocol = strings.TrimSpace(protocol)
		if strings.HasPrefix(protocol, "ticket.") {
			if ticket != "" {
				http.Error(w, "invalid ticket", http.StatusUnauthorized)
				return
			}
			ticket = strings.TrimPrefix(protocol, "ticket.")
		}
	}

	grant, ok := h.play.authorize(ticket)
	if !ok {
		http.Error(w, "game ticket required", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"angel.game"},
	})
	if err != nil {
		return
	}
	if conn.Subprotocol() != "angel.game" {
		_ = conn.Close(websocket.StatusPolicyViolation, "protocol required")
		return
	}

	conn.SetReadLimit(maxMessageSize)
	c := &client{
		hub: h, conn: conn,
		out: make(chan []byte, 64),
		id:  grant.User,
	}

	h.mu.Lock()
	for existing := range h.clients {
		if existing.id == grant.User {
			h.mu.Unlock()
			_ = conn.Close(websocket.StatusPolicyViolation, "already connected")
			return
		}
	}
	h.clients[c] = struct{}{}
	h.allowed[c] = grant.Room
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.allowed, c)
		h.mu.Unlock()
	}()

	go c.writeLoop()
	c.readLoop()
}
