package hub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 64 << 10
)

type client struct {
	hub         *Hub
	conn        *websocket.Conn
	out         chan []byte
	id          types.UserID
	room        *room.Room
	mu          sync.Mutex
	ip          string 
	fingerprint string 
}

func (c *client) roomSnapshot() *room.Room {
	c.mu.Lock(); defer c.mu.Unlock()
	return c.room
}

func (c *client) setRoom(rm *room.Room) {
	c.mu.Lock(); c.room = rm; c.mu.Unlock()
}

func (c *client) safeHandle(m wire.InboundMsg) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in command handler", "client", c.id, "panic", fmt.Sprint(r))
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "internal error"))
		}
	}()
	c.handle(m)
}

func (c *client) readLoop() {
	ctx := context.Background()
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil { break }
		m, err := wire.NewInboundMsg(data)
		if err != nil {
			logger.Error("Invalid binary frame", "client", c.id, "error", err)
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid binary frame"))
			continue
		}
		c.safeHandle(m)
	}
	c.hub.mu.Lock()
	delete(c.hub.clients, c)
	c.hub.mu.Unlock()

	rm := c.roomSnapshot()
	c.setRoom(nil)
	
	// آزادسازی قفل دستگاه هنگام قطع اتصال
	if c.fingerprint != "" {
		room.ReleaseLock(ctx, c.fingerprint, string(c.id))
	}

	if rm != nil {
		lctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = rm.Leave(lctx, c.id)
		cancel()
	}
	logger.Info("Client disconnected", "client", c.id)
}

func (c *client) writeLoop() {
	defer c.conn.CloseNow()
	ctx := context.Background()
	for {
		select {
		case <-c.hub.done: return
		case frame := <-c.out:
			// ⬇️ اصلاح شد: استفاده از writeWait به جای writeTimeout ⬇️
			wctx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(wctx, websocket.MessageBinary, frame)
			cancel()
			if err != nil { return }
		}
	}
}

func (c *client) enqueue(frame []byte) {
	select {
	case c.out <- frame:
	default:
		logger.Error("Outbound buffer full, frame dropped", "client", c.id)
	}
}
