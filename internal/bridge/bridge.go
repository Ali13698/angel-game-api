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
	done  chan struct{}
}

func New() *Bridge {
	return &Bridge{
		conns: make(map[*websocket.Conn]bool),
		done:  make(chan struct{}),
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
	if len(b.conns) == 0 { return }
	
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
	log.Printf("bridge: node connected from %s (total: %d)", r.RemoteAddr, len(b.conns))

	// 🛠️ ارسال پینگ از سمت Go به Node.js برای جلوگیری از قطعی کلودفلر
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := conn.Write(ctx, websocket.MessageText, []byte(`{"type":"ping"}`))
				cancel()
				if err != nil { return }
			case <-b.done:
				return
			}
		}
	}()

	ctx := context.Background()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("bridge: read error from %s: %v", r.RemoteAddr, err)
			break
		}
		if string(data) == `{"type":"ping"}` {
			_ = conn.Write(ctx, websocket.MessageText, []byte(`{"type":"pong"}`))
		}
	}

	b.mu.Lock()
	delete(b.conns, conn)
	b.mu.Unlock()
	log.Printf("bridge: node disconnected from %s", r.RemoteAddr)
}
