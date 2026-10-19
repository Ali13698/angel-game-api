package hub
import ("context";"time";"github.com/angel/angel/internal/core/types";"github.com/angel/angel/internal/logger";"github.com/angel/angel/internal/room";"github.com/angel/angel/internal/wire")
func (h *Hub) broadcastToRoom(rm *room.Room, frame []byte) {
	h.mu.Lock(); clients := make([]*client, 0, len(h.clients)); for c := range h.clients { clients = append(clients, c) }; h.mu.Unlock()
	for _, c := range clients { c.mu.Lock(); my := c.room; c.mu.Unlock(); if my != nil && my.ID() == rm.ID() { c.enqueue(frame) } }
	if h.br != nil { h.br.PushRaw(frame) }
}
type wireRoom struct{ h *Hub }
func (w *wireRoom) Broadcast(frame []byte) { if w.h.br != nil { w.h.br.PushRaw(frame) }; w.h.mu.Lock(); clients := make([]*client, 0, len(w.h.clients)); for c := range w.h.clients { clients = append(clients, c) }; w.h.mu.Unlock(); for _, c := range clients { c.enqueue(frame) } }
func (w *wireRoom) SendTo(id types.UserID, frame []byte) { w.h.mu.Lock(); clients := make([]*client, 0, 1); for c := range w.h.clients { if c.id == id { clients = append(clients, c) } }; w.h.mu.Unlock(); for _, c := range clients { c.enqueue(frame) } }
func (w *wireRoom) CloseSlowClient(id types.UserID, code string, reason string) { w.h.mu.Lock(); var target *client; for c := range w.h.clients { if c.id == id { target = c; break } }; w.h.mu.Unlock(); if target == nil { return }; logger.Warn("Closing slow client", "client", id, "code", code); target.enqueue(wire.ErrFrame(code, reason)); _ = target.conn.CloseNow() }
func (w *wireRoom) SignalBackpressure(code string, reason string) { w.h.mu.Lock(); clients := make([]*client, 0, len(w.h.clients)); for c := range w.h.clients { clients = append(clients, c) }; w.h.mu.Unlock(); frame := wire.ErrFrame(code, reason); for _, c := range clients { c.enqueue(frame) } }
func (h *Hub) Shutdown(ctx context.Context) error { h.once.Do(func() { close(h.done); h.mu.Lock(); for c := range h.clients { _ = c.conn.CloseNow() }; h.mu.Unlock(); logger.Info("Hub shutdown complete") }); select { case <-ctx.Done(): return ctx.Err(); case <-time.After(100 * time.Millisecond): return nil } }
