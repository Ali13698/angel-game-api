package hub

import (
	"net/http"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/room"
	"github.com/coder/websocket"
)

func (h *Hub) OpenRoom(gameID types.GameID, players []types.UserID) (*room.Room, error) {
	w := &wireRoom{h: h}
	
	rm, err := h.reg.Open(registry.Spec{
		GameID:         gameID,
		Players:        players,
		Transport:      w,
		StatsRecorder:  h.StatsRecorder,
		EventPublisher: h.EventPublisher, // ارسال کافکا به رجیستری
	})
	if err != nil {
		return nil, err
	}
	
	h.mu.Lock()
	h.rooms[rm.ID()] = rm
	h.mu.Unlock()
	logger.Info("Room opened", "room_id", rm.ID(), "game_id", gameID)
	return rm, nil
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid := types.UserID(r.URL.Query().Get("user"))
	if uid == "" {
		http.Error(w, "user required", http.StatusUnauthorized)
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		logger.Error("Hub accept failed", "error", err)
		return
	}
	conn.SetReadLimit(maxMessageSize)
	c := &client{hub: h, conn: conn, out: make(chan []byte, 64), id: uid}
	
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	
	logger.Info("Client connected", "client", uid, "addr", r.RemoteAddr)
	go c.writeLoop()
	c.readLoop()
}
