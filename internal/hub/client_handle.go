package hub

import (
	"context"
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
		st, rev := rm.GetState()
		c.enqueue(wire.OutboundMsg{Type: wire.TypeJoined, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())

	default:
		rm := c.roomSnapshot()
		if rm == nil {
			c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "join first"))
			return
		}
		cmd, err := m.Command()
		if err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, err.Error()))
			return
		}
		if cmd.RoomID == "" { cmd.RoomID = rm.ID() }
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := rm.Do(ctx, room.Request{From: c.id, Cmd: cmd}); err != nil {
			c.enqueue(errFrameFor(err))
		}
	}
}
