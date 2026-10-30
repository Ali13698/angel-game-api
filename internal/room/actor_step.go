package room

import (
	"context"
	"fmt"
	"time"

	"github.com/angel/angel/internal/core/types"
)

type actorCmdMsg struct {
	from types.UserID
	cmd  types.Command
	resp chan actorCmdResult
}

type actorCmdResult struct {
	events []types.Event
	err    error
}

func (a *RoomActor) handleCommand(ctx context.Context, msg actorCmdMsg) {
	start := time.Now()
	defer func() {
		a.stats.recordExecution(time.Since(start))
	}()

	currentRev := a.state.Version()
	if msg.cmd.Version == 0 || msg.cmd.Version != currentRev {
		msg.resp <- actorCmdResult{
			err: fmt.Errorf("room: command version mismatch: got %d, want %d", msg.cmd.Version, currentRev),
		}
		return
	}

	nextState, events, err := a.handler.ApplyCommand(a.state, msg.from, msg.cmd)
	if err != nil {
		msg.resp <- actorCmdResult{err: err}
		return
	}

	a.state = nextState
	a.broadcastEvents(events)

	msg.resp <- actorCmdResult{
		events: events,
		err:    nil,
	}
}

func (a *RoomActor) handleTimeout(ctx context.Context) {
	if a.handler.IsTerminal(a.state) {
		return
	}

	current := a.state
	cmd := types.Command{
		Type:    "system_timeout",
		Version: current.Version(),
	}

	nextState, events, err := a.handler.ApplyCommand(current, "system", cmd)
	if err != nil {
		return
	}

	a.state = nextState
	a.broadcastEvents(events)
}
