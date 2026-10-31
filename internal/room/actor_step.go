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
		if next == nil {
			next = prev
		}
		if next.Version() != prev.Version() {
			r.commit(next)
		}
		r.emit(events)
		r.reply(env, nil)
		r.checkTerminal()
		return
	}

	if r.Phase() != PhaseActive {
		r.reply(env, ErrWrongPhase)
		return
	}

	// 🛡️ Anti-Cheat Gate: Strict Version Check
	// Removed the `> 0` bypass. The client MUST send the exact current revision.
	if env.Cmd.Version < prev.Version() {
		r.reply(env, ErrStaleVersion)
		return
	}

	uid, _ := strconv.ParseInt(string(env.From), 10, 64)
	if uid > 0 {
		active, err := items.GetUserActiveItems(context.Background(), uid)
		if err == nil {
			env.Cmd.ActiveItems = items.MapToWasmItems(active)
		}
	}

	next, events, err := r.cfg.Game.ApplyCommand(prev, env.From, env.Cmd)
	if err != nil {
		r.reply(env, err)
		return
	}
	if next == nil {
		next = prev
	}
	if next.Version() != prev.Version() {
		r.commit(next)
		r.mu.Lock()
		r.lastMoveTime = time.Now()
		r.mu.Unlock()
	}

	for _, event := range events {
		r.ProcessVoiceEvent(event)
	}
	r.emit(events)
	r.reply(env, nil)
	r.checkTerminal()
}

func (r *Room) tick(now time.Time) {
	prev, _ := r.GetState()
	next, events := r.cfg.Game.Tick(prev, now)
	if next == nil {
		next = prev
	}
	if next.Version() != prev.Version() {
		r.commit(next)
	}

	if r.Phase() == PhaseActive {
		r.mu.RLock()
		lastMove := r.lastMoveTime
		r.mu.RUnlock()

		current, _ := r.GetState()
		if now.Sub(lastMove) > 10*time.Second && !r.cfg.Game.IsTerminal(current) {
			cmd := types.Command{
				Type:    "timeout",
				Version: current.Version(),
			}
			timeoutNext, timeoutEvents, err := r.cfg.Game.ApplyCommand(
				current, "server", cmd,
			)
			if err == nil {
				if timeoutNext != nil && timeoutNext.Version() != current.Version() {
					r.commit(timeoutNext)
				}
				events = append(events, timeoutEvents...)
			}
			r.mu.Lock()
			r.lastMoveTime = now
			r.mu.Unlock()
		}
	}

	for _, event := range events {
		r.ProcessVoiceEvent(event)
	}
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

func (r *Room) advancePhase(next Phase, now time.Time) {
	r.mu.Lock()
	current := r.phase
	if current == next {
		r.mu.Unlock()
		return
	}

	allowed := false
	switch current {
	case PhaseLobby:
		allowed = next == PhaseCountdown || next == PhaseAborted
	case PhaseCountdown:
		allowed = next == PhaseActive || next == PhaseAborted
	case PhaseActive:
		allowed = next == PhaseFinished || next == PhaseAborted
	case PhaseFinished:
		allowed = next == PhaseAborted
	case PhaseAborted:
		allowed = false
	}
	if !allowed {
		r.mu.Unlock()
		return
	}

	r.phase = next
	r.phaseAt = now
	if next == PhaseActive {
		r.lastMoveTime = now
	}
	r.mu.Unlock()

	logger.Info(
		"Room phase changed",
		"room_id", r.cfg.ID,
		"phase", next.String(),
	)
	if next == PhaseActive {
		r.emit([]types.Event{
			{Type: "game.started", Audience: types.ToRoom},
		})
	}
}

func (r *Room) checkTerminal() {
	if r.Phase() != PhaseActive {
		return
	}
	state, _ := r.GetState()
	if r.cfg.Game.IsTerminal(state) {
		r.handleGameEnd(state)
	}
}

func (r *Room) commit(next types.State) {
	r.mu.Lock()
	r.cur = next
	r.mu.Unlock()
}
