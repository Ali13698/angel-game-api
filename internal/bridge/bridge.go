package bridge

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Bridge struct {
	mu    sync.RWMutex
	conns map[*websocket.Conn]bool
}

func New() *Bridge {
	return &Bridge{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (b *Bridge) Connected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.conns) > 0
}

func (b *Bridge) PushRaw(frame []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for conn := range b.conns {
		go func(c *websocket.Conn) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = c.Write(ctx, websocket.MessageText, frame)
		}(conn)
	}
}

func (b *Bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("bridge: accept: %v", err)
		return
	}
	conn.SetReadLimit(1 << 20)

	b.mu.Lock()
	b.conns[conn] = true
	b.mu.Unlock()

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := conn.Write(ctx, websocket.MessageText, []byte(`{"type":"ping"}`))
			cancel()
			if err != nil {
				return
			}
		}
	}()

	ctx := context.Background()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			break
		}
		if string(data) == `{"type":"ping"}` {
			_ = conn.Write(ctx, websocket.MessageText, []byte(`{"type":"pong"}`))
		}
	}

	b.mu.Lock()
	delete(b.conns, conn)
	b.mu.Unlock()
}
