package hub
import ("context";"encoding/json";"time";"github.com/angel/angel/internal/core/types";"github.com/angel/angel/internal/logger";"github.com/angel/angel/internal/room";"github.com/angel/angel/internal/wire";"github.com/coder/websocket")
func (c *client) readLoop() {
	ctx := context.Background(); for { _, data, err := c.conn.Read(ctx); if err != nil { break }; m, err := wire.NewInboundMsg(data); if err != nil { logger.Error("Invalid binary frame", "client", c.id); c.enqueue(wire.ErrFrame(wire.CodeBadFrame, "invalid binary frame")); continue }; c.handle(m) }
	c.hub.mu.Lock(); delete(c.hub.clients, c); c.hub.mu.Unlock(); c.mu.Lock(); rm := c.room; c.room = nil; c.mu.Unlock(); if rm != nil { _ = rm.Leave(ctx, c.id) }; logger.Info("Client disconnected", "client", c.id)
}
func (c *client) writeLoop() { defer c.conn.CloseNow(); ctx := context.Background(); for { select { case <-c.hub.done: return; case frame := <-c.out: wctx, cancel := context.WithTimeout(ctx, writeWait); err := c.conn.Write(wctx, websocket.MessageBinary, frame); cancel(); if err != nil { logger.Error("Failed to write to client", "client", c.id, "error", err); return } } } }
func (c *client) enqueue(frame []byte) { select { case c.out <- frame: default: } }
func (c *client) handle(m wire.InboundMsg) {
	switch m.Type {
	case wire.TypePing: c.enqueue(wire.OutboundMsg{Type: wire.TypePong}.Bytes())
	case wire.TypeJoin:
		c.hub.mu.Lock(); rm, ok := c.hub.rooms[types.RoomID(m.RoomID)]; c.hub.mu.Unlock(); if !ok { c.enqueue(wire.ErrFrame(wire.CodeNoRoom, "no such room")); return }
		c.mu.Lock(); c.room = rm; c.mu.Unlock()
		{ jctx, jcancel := context.WithTimeout(context.Background(), 5*time.Second); if err := rm.Join(jctx, c.id); err != nil { jcancel(); c.enqueue(errFrameFor(err)); return }; jcancel() }
		st, rev := rm.GetState(); c.enqueue(wire.OutboundMsg{Type: wire.TypeJoined, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())
	case wire.TypeLeave: c.mu.Lock(); rm := c.room; c.room = nil; c.mu.Unlock(); if rm != nil { _ = rm.Leave(context.Background(), c.id) }
	case wire.TypeChat:
		c.mu.Lock(); rm := c.room; c.mu.Unlock(); if rm == nil { return }
		var p struct{ Text string `json:"text"` }; if err := json.Unmarshal(m.Payload, &p); err != nil { return }; if len(p.Text) > 500 { p.Text = p.Text[:500] }; if p.Text == "" { return }
		c.hub.broadcastToRoom(rm, wire.OutboundMsg{Type: wire.TypeChat, RoomID: string(rm.ID()), From: string(c.id), Payload: map[string]any{"text": p.Text, "user_id": string(c.id)}}.Bytes())
	case wire.TypeEmoji:
		c.mu.Lock(); rm := c.room; c.mu.Unlock(); if rm == nil { return }
		var p struct{ Emoji string `json:"emoji"` }; if err := json.Unmarshal(m.Payload, &p); err != nil { return }; if len(p.Emoji) > 8 { p.Emoji = p.Emoji[:8] }; if p.Emoji == "" { return }
		c.hub.broadcastToRoom(rm, wire.OutboundMsg{Type: wire.TypeEmoji, RoomID: string(rm.ID()), From: string(c.id), Payload: map[string]any{"emoji": p.Emoji, "user_id": string(c.id)}}.Bytes())
	case wire.TypeState:
		c.mu.Lock(); rm := c.room; c.mu.Unlock(); if rm == nil { c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "not joined")); return }
		st, rev := rm.GetState(); c.enqueue(wire.OutboundMsg{Type: wire.TypeState, RoomID: string(rm.ID()), GameID: string(rm.GameID()), Rev: rev, State: st}.Bytes())
	default:
		c.mu.Lock(); rm := c.room; c.mu.Unlock(); if rm == nil { c.enqueue(wire.ErrFrame(wire.CodeNotJoined, "join first")); return }
		cmd, err := m.Command(); if err != nil { c.enqueue(wire.ErrFrame(wire.CodeBadFrame, err.Error())); return }; if cmd.RoomID == "" { cmd.RoomID = rm.ID() }
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel(); if err := rm.Do(ctx, room.Request{From: c.id, Cmd: cmd}); err != nil { c.enqueue(errFrameFor(err)) }
	}
}
