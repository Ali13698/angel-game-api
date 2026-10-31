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
		requested := types.RoomID(m.RoomID)

		c.hub.mu.Lock()
		allowed, authorized := c.hub.allowed[c]
		rm, exists := c.hub.rooms[requested]
		c.hub.mu.Unlock()

		if !authorized || requested != allowed {
			c.enqueue(wire.ErrFrame(wire.CodeRejected, "room permission required"))
			return
		}
		if !exists {
			c.enqueue(wire.ErrFrame(wire.CodeNoRoom, "no such room"))
			return
		}

		if prev := c.roomSnapshot(); prev == rm {
			c.enqueue(joinFrame(rm))
			return
		} else if prev != nil {
			c.enqueue(wire.ErrFrame(wire.CodeRejected, "already in another room"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := rm.Join(ctx, c.id)
		cancel()
		if err != nil {
			c.enqueue(errFrameFor(err))
			return
		}

		c.setRoom(rm)
		c.enqueue(joinFrame(rm))

	default:
		rm := c.roomSnapshot()
		if rm == nil {
			c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "join first"))
			return
		}
		if m.RoomID != "" && types.RoomID(m.RoomID) != rm.ID() {
			c.enqueue(wire.ErrFrame(wire.CodeRejected, "room mismatch"))
			return
		}

		cmd, err := m.Command()
		if err != nil {
			c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid command"))
			return
		}
		cmd.RoomID = rm.ID()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := rm.Do(ctx, room.Request{From: c.id, Cmd: cmd}); err != nil {
			c.enqueue(errFrameFor(err))
		}
	}
}

func joinFrame(rm *room.Room) []byte {
	st, rev := rm.GetState()
	var state []byte
	if st != nil {
		state = wire.EncodeValue(st)
	}
	return wire.OutboundMsg{
		Type:   wire.TypeJoined,
		RoomID: string(rm.ID()),
		GameID: string(rm.GameID()),
		Rev:    rev,
		State:  state,
	}.Bytes()
}
