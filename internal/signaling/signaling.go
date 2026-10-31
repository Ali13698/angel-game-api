package signaling

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[string]*websocket.Conn)}
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
	userID := r.URL.Query().Get("user_id")
	if roomID == "" || userID == "" {
		http.Error(w, "room_id and user_id required", http.StatusBadRequest)
		return
	}
	if !originAllowed(r) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
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
		typ, data, err := conn.Read(context.Background())
		if err != nil {
			return
		}
		if typ != websocket.MessageBinary {
			continue
		}
		var msg wire.ClientMsg
		if err := proto.Unmarshal(data, &msg); err != nil {
			continue
		}
		body, err := wire.DecodeMap(msg.Payload)
		if err != nil {
			continue
		}
		target, _ := body["target"].(string)
		if target == "" {
			continue
		}

		out := wire.OutboundMsg{Type: "signal", From: userID, Payload: msg.Payload}.Bytes()
		if out == nil {
			continue
		}

		h.mu.RLock()
		targetConn := h.rooms[roomID][target]
		h.mu.RUnlock()

		if targetConn != nil {
			_ = targetConn.Write(context.Background(), websocket.MessageBinary, out)
		}
	}
}
