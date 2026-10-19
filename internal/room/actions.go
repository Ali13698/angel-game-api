package room
import ("context";"time";"github.com/angel/angel/internal/core/types")
func (r *Room) Join(ctx context.Context, id types.UserID) error {
	r.mu.Lock(); prev := r.cur; r.mu.Unlock(); if prev == nil { return ErrClosed }
	next, events := r.cfg.Game.OnPlayerJoined(prev, id); if next != nil && next.Version() != prev.Version() { r.commit(next) }
	r.emit(events); r.mu.Lock(); r.joinedCount++; jc := r.joinedCount; ph := r.phase; r.mu.Unlock()
	if jc >= 2 && ph == PhaseLobby { r.advancePhase(PhaseCountdown, time.Now()) }; return nil
}
func (r *Room) Leave(ctx context.Context, id types.UserID) error {
	reply := make(chan error, 1); env := Envelope{From: id, Leave: true, ReplyTo: reply}
	select { case r.priority <- env: default:
		waitCtx, cancel := context.WithTimeout(ctx, PriorityLaneWaitTimeout); defer cancel()
		select { case r.priority <- env: case <-r.quit: return ErrClosed; case <-waitCtx.Done(): if ctx.Err() != nil { return ctx.Err() }; return ErrPriorityFull }
	}
	select { case err := <-reply: return err; case <-r.quit: return ErrClosed; case <-ctx.Done(): return ctx.Err() }
}
