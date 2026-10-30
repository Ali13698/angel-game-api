package room

import (
	"context"
	"time"

	"github.com/angel/angel/internal/wire"
)

func (r *RoomActor) step(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case env := <-r.cmdChan:
		prev := r.state
		
		// Phase 1 strict revision validation
		if env.Cmd.Version == 0 || env.Cmd.Version != prev.Version() {
			env.ErrChan <- ErrRevisionMismatch
			return
		}

		newState, events, err := r.cfg.Game.ApplyCommand(ctx, r.cfg.GameID, prev.Bytes(), env.Cmd)
		if err != nil {
			env.ErrChan <- err
			r.metrics.IncError()
			return
		}

		st, err := ParseState(newState)
		if err != nil {
			env.ErrChan <- err
			return
		}

		r.state = st
		r.metrics.IncStep()
		env.ErrChan <- nil

		r.broadcastEvents(events)
		r.checkTerminal()

	case <-r.timeoutChan:
		// Check if room is idle
		if time.Since(r.lastActive) > r.cfg.IdleTimeout {
			r.phase = PhaseTerminated
		}

	case req := <-r.syncChan:
		req.RespChan <- r.snapshot()

	case <-ctx.Done():
		r.phase = PhaseTerminated
	}
}

func (r *RoomActor) checkTerminal() {
	if r.state.IsTerminal() {
		r.phase = PhaseTerminated
	}
}

func (r *RoomActor) tickTimeout() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.phase == PhaseTerminated {
		return
	}

	if time.Since(r.lastActive) > r.cfg.TurnTimeout {
		// Auto-forfeit or timeout command with up-to-date revision
		current, _ := r.GetState()
		timeoutCmd := wire.Command{
			Type:    "TIMEOUT",
			Player:  r.currentTurnPlayer(),
			Version: current.Version(),
		}
		
		newState, events, err := r.cfg.Game.ApplyCommand(context.Background(), r.cfg.GameID, r.state.Bytes(), timeoutCmd)
		if err == nil {
			if st, parseErr := ParseState(newState); parseErr == nil {
				r.state = st
				r.lastActive = time.Now()
				r.broadcastEvents(events)
				r.checkTerminal()
			}
		}
	}
}
