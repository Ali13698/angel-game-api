package hub

import (
	"net/http"
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
	done           chan struct{}
	once           sync.Once
}

func New(reg *registry.Registry, pub room.EventPublisher) *Hub {
	return &Hub{
		reg: reg, EventPublisher: pub,
		clients: make(map[*client]struct{}),
		rooms:   make(map[types.RoomID]*room.Room),
		done:    make(chan struct{}),
	}
}

func (h *Hub) SetBridge(br *bridge.Bridge) { h.br = br }

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	w := &wireRoom{h: h}
	rm, err := h.reg.Open(registry.Spec{
		GameID: gameID, Players: players, Transport: w, StatsRecorder: h.StatsRecorder, EventPublisher: h.EventPublisher,
	})
	if err != nil { return nil, err }

	h.mu.Lock()
	h.rooms[rm.ID()] = rm
	h.mu.Unlock()
	logger.Info("Room opened", "room_id", rm.ID(), "game_id", gameID)

	hasBot := false
	for _, p := range players {
		if strings.HasPrefix(string(p), "bot") {
			hasBot = true
			break
		}
	}
	if hasBot {
		// ✅ تغییر شناسه به "bot" برای هماهنگی با فرانت‌اند
		go h.runBot(rm, types.UserID("bot"), gameID, 3*time.Second)
	}

	return rm, nil
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := types.UserID(r.URL.Query().Get("user"))
	if uid == "" { http.Error(w, "user required", http.StatusUnauthorized); return }
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil { return }
	conn.SetReadLimit(maxMessageSize)
	c := &client{hub: h, conn: conn, out: make(chan []byte, 64), id: uid}
	
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	
	go c.writeLoop()
	c.readLoop()
}
