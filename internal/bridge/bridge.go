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

	ctx, cancel := context.WithCancel(r.Context()) // 🔴 استفاده از Context برای بستن حلقه
	defer cancel()

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				writeCtx, writeCancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(writeCtx, websocket.MessageText, []byte(`{"type":"ping"}`))
				writeCancel()
				if err != nil {
					return
				}
			}
		}
	}()

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
