package hub

import (
	"context"
	"fmt"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
)

func (c *client) handle(m wire.InboundMsg) {
	switch m.Type {
	case wire.TypePing:
		c.enqueue(wire.OutboundMsg{Type: wire.TypePong}.Bytes())

	case wire.TypeJoin:
		// ۱. استخراج Fingerprint از پیلود (اگر فرانت‌اند فرستاده باشد)
		// در غیر این صورت از IP به عنوان شناسه دستگاه استفاده می‌کنیم
		fingerprint := c.ip // fallback
		if fp, ok := m.Payload["fingerprint"].(string); ok && fp != "" {
			fingerprint = fp
		}

		// ۲. تلاش برای گرفتن قفل در Redis
		if err := room.AcquireLock(context.Background(), fingerprint, string(c.id)); err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeRejected, err.Error()))
			return
		}

		// ۳. ادامه منطق ورود به اتاق (مثل قبل)
		c.hub.mu.Lock()
		rm, ok := c.hub.rooms[types.RoomID(m.RoomID)]
		c.hub.mu.Unlock()
		if !ok {
			c.enqueue(wire.ErrFrame(wire.CodeNoRoom, "no such room"))
			return
		}
		prev := c.roomSnapshot()
		if prev == rm {
			st, rev := rm.GetState()
			c.enqueue(wire.OutboundMsg{Type: wire.TypeJoined, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())
			return
		}
		if prev != nil {
			lctx, lcancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = prev.Leave(lctx, c.id)
			lcancel()
			c.setRoom(nil)
		}
		
		jctx, jcancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := rm.Join(jctx, c.id)
		jcancel()
		if err != nil {
			c.enqueue(errFrameFor(err))
			return
		}
		
		c.setRoom(rm)
		c.fingerprint = fingerprint // ذخیره برای آزادسازی در زمان خروج
		
		st, rev := rm.GetState()
		c.enqueue(wire.OutboundMsg{Type: wire.TypeJoined, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())

	default:
		// ... (منطق دستورات بازی مثل roll و move)
	}
}
