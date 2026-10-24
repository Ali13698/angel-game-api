package room

import (
	"context"
	"strconv"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/items"
	"github.com/angel/angel/internal/logger"
)

func (r *Room) step(env Envelope) {
	prev, _ := r.GetState()
	if env.Leave {
		next, events := r.cfg.Game.OnPlayerLeft(prev, env.From)
		if next == nil { next = prev }
		if next.Version() != prev.Version() { r.commit(next) }
		r.emit(events)
		r.reply(env, nil)
		r.checkTerminal()
		return
	}

	if r.Phase() != PhaseActive {
		logger.Warn("⚠️ [ROOM] Rejected: wrong phase", "room_id", r.cfg.ID, "phase", r.Phase())
		r.reply(env, ErrWrongPhase)
		return
	}

	if env.Cmd.Version > 0 && env.Cmd.Version < prev.Version() {
		logger.Warn("⚠️ [ROOM] Rejected: stale version", "room_id", r.cfg.ID)
		r.reply(env, ErrStaleVersion)
		return
	}

	uid, _ := strconv.ParseInt(string(env.From), 10, 64)
	if uid > 0 {
		actItems, err := items.GetUserActiveItems(context.Background(), uid)
		if err == nil {
			env.Cmd.ActiveItems = items.MapToWasmItems(actItems)
		}
	}

	logger.Info("🧠 [ROOM] Calling WASM ApplyCommand", "room_id", r.cfg.ID, "type", env.Cmd.Type, "from", env.From)
	next, events, err := r.cfg.Game.ApplyCommand(prev, env.From, env.Cmd)
	if err != nil {
		logger.Error("❌ [ROOM] WASM returned error", "room_id", r.cfg.ID, "error", err)
		r.reply(env, err)
		return
	}

	if next == nil { next = prev }
	if next.Version() != prev.Version() {
		r.commit(next)
	}

	for _, ev := range events {
		logger.Info("🎉 [ROOM] WASM generated event", "room_id", r.cfg.ID, "event", ev.Type)
		r.ProcessVoiceEvent(ev)
	}

	r.emit(events)
	r.reply(env, nil)
	r.checkTerminal()
}

func (r *Room) tick(now time.Time) {
	prev, _ := r.GetState()
	next, events := r.cfg.Game.Tick(prev, now)
	if next == nil { next = prev }
	if next.Version() != prev.Version() { r.commit(next) }
	for _, ev := range events { r.ProcessVoiceEvent(ev) }
	r.emit(events)

	if r.Phase() == PhaseCountdown {
		r.mu.RLock()
		started := r.phaseAt
		r.mu.RUnlock()
		if now.Sub(started) >= r.countdown {
			r.advancePhase(PhaseActive, now)
		}
	}
	r.checkTerminal()
}

// ✅ تابع اصلی advancePhase (تنها نسخه در پروژه)
func (r *Room) advancePhase(next Phase, now time.Time) {
	r.mu.Lock()
	cur := r.phase
	if cur == next { r.mu.Unlock(); return }
	allowed := false
	switch cur {
	case PhaseLobby: allowed = next == PhaseCountdown || next == PhaseAborted
	case PhaseCountdown: allowed = next == PhaseActive || next == PhaseAborted
	case PhaseActive: allowed = next == PhaseFinished || next == PhaseAborted
	case PhaseFinished: allowed = next == PhaseAborted
	case PhaseAborted: allowed = false
	}
	if !allowed { r.mu.Unlock(); return }
	r.phase = next
	r.phaseAt = now
	r.mu.Unlock()
	logger.Info("🔄 [ROOM] Phase changed", "room_id", r.cfg.ID, "phase", next.String())

	// ✅ ارسال رویداد شروع بازی به مرورگر
	if next == PhaseActive {
		r.emit([]types.Event{{Type: "game.started", Audience: types.ToRoom}})
	}
}

func (r *Room) checkTerminal() {
	if r.Phase() != PhaseActive { return }
	st, _ := r.GetState()
	if r.cfg.Game.IsTerminal(st) {
		r.handleGameEnd(st)
	}
}

func (r *Room) commit(next types.State) {
	r.mu.Lock()
	r.cur = next
	r.mu.Unlock()
}
