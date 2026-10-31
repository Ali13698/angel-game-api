package bridge

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
)

const (
	writeTimeout = 5 * time.Second
	pingInterval = 15 * time.Second
	readLimit    = 1 << 20
)

type Bridge struct {
	mu    sync.RWMutex
	conns map[*websocket.Conn]bool
}

func New() *Bridge {
	return &Bridge{conns: make(map[*websocket.Conn]bool)}
}

func (b *Bridge) Connected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.conns) > 0
}

func (b *Bridge) PushRaw(frame []byte) {
	b.mu.RLock()
	targets := make([]*websocket.Conn, 0, len(b.conns))
	for conn := range b.conns {
		targets = append(targets, conn)
	}
	b.mu.RUnlock()

	for _, conn := range targets {
		go func(c *websocket.Conn) {
			ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
			defer cancel()
			if err := c.Write(ctx, websocket.MessageBinary, frame); err != nil {
				b.drop(c)
			}
		}(conn)
	}
}

func (b *Bridge) drop(c *websocket.Conn) {
	b.mu.Lock()
	delete(b.conns, c)
	b.mu.Unlock()
	_ = c.Close(websocket.StatusInternalError, "write failed")
}

func (b *Bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isWebSocketUpgrade(r) {
		http.Error(w, "websocket upgrade required", http.StatusUpgradeRequired)
		return
	}
	if !originAllowed(r) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		log.Printf("bridge: accept: %v", err)
		return
	}
	conn.SetReadLimit(readLimit)
	log.Printf("bridge: bot connected from %s", r.RemoteAddr)

	b.mu.Lock()
	b.conns[conn] = true
	b.mu.Unlock()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	defer func() {
		b.mu.Lock()
		delete(b.conns, conn)
		b.mu.Unlock()
		_ = conn.CloseNow()
		log.Printf("bridge: bot disconnected from %s", r.RemoteAddr)
	}()

	go func() {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				writeCtx, writeCancel := context.WithTimeout(context.Background(), writeTimeout)
				err := conn.Write(writeCtx, websocket.MessageBinary, pingFrame())
				writeCancel()
				if err != nil {
					cancel()
					return
				}
			}
		}
	}()

	for {
		typ, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		if typ != websocket.MessageBinary {
			continue
		}
		if msg, perr := wire.NewInboundMsg(data); perr == nil && msg.Type == wire.TypePing {
			writeCtx, writeCancel := context.WithTimeout(context.Background(), writeTimeout)
			_ = conn.Write(writeCtx, websocket.MessageBinary, wire.OutboundMsg{Type: wire.TypePong}.Bytes())
			writeCancel()
		}
	}
}

func pingFrame() []byte {
	return wire.OutboundMsg{Type: wire.TypePing}.Bytes()
}
