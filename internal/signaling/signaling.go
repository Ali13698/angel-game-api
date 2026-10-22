package signaling

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

// Hub مدیریت سیگنالینگ WebRTC برای تماس صوتی
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]*websocket.Conn // roomID -> userID -> conn
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[string]*websocket.Conn)}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	userID := r.URL.Query().Get("user_id")
	if roomID == "" || userID == "" {
		http.Error(w, "room_id and user_id required", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		return
	}

	h.mu.Lock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[string]*websocket.Conn)
	}
	h.rooms[roomID][userID] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.rooms[roomID], userID)
		h.mu.Unlock()
		conn.CloseNow()
	}()

	for {
		_, data, err := conn.Read(context.Background())
		if err != nil {
			return
		}
		
		var msg struct {
			Target string `json:"target"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		
		h.mu.RLock()
		targetConn := h.rooms[roomID][msg.Target]
		h.mu.RUnlock()
		
		if targetConn != nil {
			_ = targetConn.Write(context.Background(), websocket.MessageText, data)
		}
	}
}
