package chat

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/coder/websocket"
)

const chatMaxFrame = 64 << 10

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*websocket.Conn]bool)}
}

func originAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	if origin == "null" {
		return false
	}
	host := strings.ToLower(r.Host)
	for _, suffix := range []string{"nomi98.com", "nomi98.ir"} {
		if strings.HasSuffix(host, suffix) {
			return strings.Contains(origin, suffix)
		}
	}
	return strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "http://127.0.0.1")
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "room_id required", http.StatusBadRequest)
		return
	}
	if !originAllowed(r) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return
	}
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		http.Error(w, "websocket upgrade required", http.StatusUpgradeRequired)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		return
	}
	conn.SetReadLimit(chatMaxFrame)

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
		typ, data, err := conn.Read(context.Background())
		if err != nil {
			return
		}
		if typ != websocket.MessageBinary || len(data) == 0 {
			continue
		}

		h.mu.RLock()
		for c := range h.rooms[roomID] {
			if c != conn {
				_ = c.Write(context.Background(), websocket.MessageBinary, data)
			}
		}
		h.mu.RUnlock()
	}
}
