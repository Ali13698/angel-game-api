package chat

import (
	"context"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

// Hub مدیریت چت متنی را بر عهده دارد و هیچ دخلی به Room Actor ندارد
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*websocket.Conn]bool)}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "room_id required", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		return
	}

	h.mu.Lock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[roomID][conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.rooms[roomID], conn)
		h.mu.Unlock()
		conn.CloseNow()
	}()

	for {
		_, data, err := conn.Read(context.Background())
		if err != nil {
			return
		}
		
		// انتشار پیام به بقیه اعضای اتاق
		h.mu.RLock()
		for c := range h.rooms[roomID] {
			if c != conn {
				_ = c.Write(context.Background(), websocket.MessageText, data)
			}
		}
		h.mu.RUnlock()
	}
}
