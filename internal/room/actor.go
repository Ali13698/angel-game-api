package room
import ("context";"sync";"time";"github.com/angel/angel/internal/core/types";"github.com/angel/angel/internal/logger")
func (r *Room) enqueueTick() { select { case r.tickCh <- struct{}{}: default: } }
func (r *Room) Run(ctx context.Context) error {
	defer close(r.done); go r.tickerLoop(ctx)
	var fanoutWG sync.WaitGroup; fanoutWG.Add(1)
	go func() { defer fanoutWG.Done(); defer close(r.fanoutDone); r.fanoutLoop(ctx) }()
	r.advancePhase(PhaseLobby, time.Now()); err := r.actorLoop(ctx)
	r.shutdownOnce.Do(func() { close(r.shutdown) }); fanoutWG.Wait(); return err
}
func (r *Room) tickerLoop(ctx context.Context) { t := time.NewTicker(time.Second); defer t.Stop(); for { select { case <-ctx.Done(): return; case <-r.quit: return; case <-r.shutdown: return; case <-t.C: r.enqueueTick() } } }
func (r *Room) actorLoop(ctx context.Context) error { for { select { case <-ctx.Done(): r.advancePhase(PhaseAborted, time.Now()); return ctx.Err(); case <-r.quit: r.advancePhase(PhaseAborted, time.Now()); return nil; case env := <-r.priority: r.step(env); case env := <-r.inbox: r.step(env); case <-r.tickCh: r.tick(time.Now()) } } }
func (r *Room) step(env Envelope) {
	prev, _ := r.GetState(); if env.Leave { next, events := r.cfg.Game.OnPlayerLeft(prev, env.From); if next == nil { next = prev }; if next.Version() != prev.Version() { r.commit(next) }; r.emit(events); r.reply(env, nil); r.checkTerminal(); return }
	if r.Phase() != PhaseActive { r.reply(env, ErrWrongPhase); return }; if env.Cmd.Version > 0 && env.Cmd.Version < prev.Version() { r.reply(env, ErrStaleVersion); return }
	next, events, err := r.cfg.Game.ApplyCommand(prev, env.From, env.Cmd); if err != nil { logger.Error("Failed to apply command", "room_id", r.cfg.ID, "error", err); r.reply(env, err); return }
	if next == nil { next = prev }; if next.Version() != prev.Version() { r.commit(next) }; r.emit(events); r.reply(env, nil); r.checkTerminal()
}
func (r *Room) tick(now time.Time) {
	prev, _ := r.GetState(); next, events := r.cfg.Game.Tick(prev, now); if next == nil { next = prev }; if next.Version() != prev.Version() { r.commit(next) }; r.emit(events)
	if r.Phase() == PhaseCountdown { r.mu.RLock(); started := r.phaseAt; r.mu.RUnlock(); if now.Sub(started) >= r.countdown { r.advancePhase(PhaseActive, now) } }; r.checkTerminal()
}
func (r *Room) advancePhase(next Phase, now time.Time) {
	r.mu.Lock(); cur := r.phase; if cur == next { r.mu.Unlock(); return }; allowed := false
	switch cur { case PhaseLobby: allowed = next == PhaseCountdown || next == PhaseAborted; case PhaseCountdown: allowed = next == PhaseActive || next == PhaseAborted; case PhaseActive: allowed = next == PhaseFinished || next == PhaseAborted; case PhaseFinished: allowed = next == PhaseAborted; case PhaseAborted: allowed = false }
	if !allowed { r.mu.Unlock(); return }; r.phase = next; r.phaseAt = now; r.mu.Unlock(); logger.Info("Phase changed", "room_id", r.cfg.ID, "phase", next.String())
}
func (r *Room) checkTerminal() {
	if r.Phase() != PhaseActive { return }; st, _ := r.GetState()
	if r.cfg.Game.IsTerminal(st) { r.advancePhase(PhaseFinished, time.Now()); if r.cfg.StatsRecorder != nil { res := r.cfg.Game.GetResult(st); go func() { ctx := context.Background(); if res.IsDraw { _ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), false, true); _ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, true) } else { _ = r.cfg.StatsRecorder.RecordMatch(ctx, res.WinnerID, string(r.cfg.GameID), true, false); _ = r.cfg.StatsRecorder.RecordMatch(ctx, res.LoserID, string(r.cfg.GameID), false, false) } }() } }
}
func (r *Room) commit(next types.State) { r.mu.Lock(); r.cur = next; r.mu.Unlock() }
