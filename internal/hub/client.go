package hub

import (
	"context"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
	"github.com/coder/websocket"
)

func (c *client) readLoop() {
	ctx := context.Background()
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil { break }
		m, err := wire.NewInboundMsg(data)
		if err != nil {
			logger.Error("Invalid binary frame", "client", c.id)
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid binary frame"))
			continue
		}
		c.safeHandle(m) // ✅ استفاده از safeHandle برای جلوگیری از کرش
	}
	c.hub.mu.Lock()
	delete(c.hub.clients, c)
	c.hub.mu.Unlock()
	
	c.mu.Lock()
	rm := c.room
	c.room = nil
	c.mu.Unlock()
	
	if rm != nil { _ = rm.Leave(ctx, c.id) }
	logger.Info("Client disconnected", "client", c.id)
}

// ✅ اضافه شدن Recover برای جلوگیری از نشتی حافظه
func (c *client) safeHandle(m wire.InboundMsg) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in handler", "client", c.id, "panic", r)
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "internal error"))
		}
	}()
	c.handle(m)
}

func (c *client) writeLoop() {
	defer c.conn.CloseNow()
	ctx := context.Background()
	for {
		select {
		case <-c.hub.done:
			return
		case frame := <-c.out:
			wctx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(wctx, websocket.MessageBinary, frame)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

func (c *client) enqueue(frame []byte) {
	select {
	case c.out <- frame:
	default:
		logger.Warn("Client buffer full, dropping frame", "client", c.id) // ✅ لاگ کردن ریزش
	}
}

func (c *client) handle(m wire.InboundMsg) {
	switch m.Type {
	case wire.TypePing:
		c.enqueue(wire.OutboundMsg{Type: wire.TypePong}.Bytes())

	case wire.TypeJoin:
		c.hub.mu.Lock()
		rm, ok := c.hub.rooms[types.RoomID(m.RoomID)]
		c.hub.mu.Unlock()
		if !ok {
			c.enqueue(wire.ErrFrame(wire.CodeNoRoom, "no such room"))
			return
		}
		
		// ✅ خروج از اتاق قبلی
		c.mu.Lock()
		oldRm := c.room
		c.mu.Unlock()
		if oldRm != nil && oldRm != rm {
			_ = oldRm.Leave(context.Background(), c.id)
		}

		jctx, jcancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := rm.Join(jctx, c.id); err != nil {
			jcancel()
			c.enqueue(errFrameFor(err))
			return
		}
		jcancel()
		
		// ✅ ست کردن اتاق فقط بعد از Join موفق
		c.mu.Lock()
		c.room = rm
		c.mu.Unlock()
		
		st, rev := rm.GetState()
		c.enqueue(wire.OutboundMsg{Type: wire.TypeJoined, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())

	default:
		c.mu.Lock()
		rm := c.room
		c.mu.Unlock()
		
		if rm == nil {
			c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "join first"))
			return
		}
		
		cmd, err := m.Command()
		if err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, err.Error()))
			return
		}
		if cmd.RoomID == "" {
			cmd.RoomID = rm.ID()
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := rm.Do(ctx, room.Request{From: c.id, Cmd: cmd}); err != nil {
			c.enqueue(errFrameFor(err))
		}
	}
}
