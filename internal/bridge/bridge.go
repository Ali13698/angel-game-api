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
	mu   sync.Mutex
	conn *websocket.Conn
	up   bool
	out  chan []byte
	done chan struct{}
}

func New() *Bridge {
	return &Bridge{
		out:  make(chan []byte, 512),
		done: make(chan struct{}),
	}
}

func (b *Bridge) Connected() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.up
}

func (b *Bridge) PushRaw(frame []byte) {
	b.mu.Lock()
	ok := b.up
	b.mu.Unlock()
	if !ok { return }
	select {
	case b.out <- frame:
	default:
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
	if b.conn != nil { _ = b.conn.CloseNow() }
	b.conn = conn
	b.up = true
	b.mu.Unlock()
	log.Printf("bridge: node connected from %s", r.RemoteAddr)

	go func() {
		ctx := context.Background()
		for {
			select {
			case <-b.done:
				return
			case data := <-b.out:
				wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err := conn.Write(wctx, websocket.MessageText, data)
				cancel()
				if err != nil { 
					log.Printf("bridge: write error: %v", err)
					return 
				}
			}
		}
	}()

	ctx := context.Background()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("bridge: read error: %v", err) 
			break
		}
		
		if string(data) == `{"type":"ping"}` {
			if err := conn.Write(ctx, websocket.MessageText, []byte(`{"type":"pong"}`)); err != nil {
				log.Printf("bridge: pong write error: %v", err)
				break
			}
		}
	}

	b.mu.Lock()
	if b.conn == conn {
		b.conn = nil
		b.up = false
	}
	b.mu.Unlock()
	log.Printf("bridge: node disconnected")
}
