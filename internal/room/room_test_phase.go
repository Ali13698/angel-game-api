package room

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

func TestRoom_PhaseTransitions_Allowed(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	if r.Phase() != PhaseLobby { t.Fatalf("initial = %v", r.Phase()) }
	now := time.Now()
	r.advancePhase(PhaseCountdown, now)
	r.advancePhase(PhaseActive, now)
	r.advancePhase(PhaseFinished, now)
	if r.Phase() != PhaseFinished { t.Fatalf("final = %v", r.Phase()) }
}

func TestRoom_PhaseTransitions_Disallowed(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	now := time.Now()
	r.advancePhase(PhaseActive, now)
	if r.Phase() != PhaseLobby { t.Fatalf("skip COUNTDOWN should reject, got %v", r.Phase()) }
	r.advancePhase(PhaseAborted, now)
	r.advancePhase(PhaseCountdown, now)
	if r.Phase() != PhaseAborted { t.Fatalf("ABORTED should be terminal, got %v", r.Phase()) }
}

func TestRoom_CheckTerminal(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	now := time.Now()
	r.advancePhase(PhaseCountdown, now)
	r.advancePhase(PhaseActive, now)
	r.cfg.Game.(*testGame).terminal.Store(true)
	r.checkTerminal()
	if r.Phase() != PhaseFinished { t.Fatalf("phase = %v", r.Phase()) }
}

func TestRoom_Inbox_Overflow(t *testing.T) {
	r := newTestRoom(t, 2, nullTransport{}, nil)
	for i := 0; i < 2; i++ { r.inbox <- Envelope{From: "x"} }
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "noop"}}); !errors.Is(err, ErrBusy) {
		t.Fatalf("err = %v, want ErrBusy", err)
	}
}

func TestRoom_TickSignal_Coalesces(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	r.enqueueTick()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ { r.enqueueTick() }
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second): t.Fatal("enqueueTick blocked")
	}
}

func TestRoom_TickSignal_BufferIsOne(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	r.enqueueTick()
	r.enqueueTick()
	r.enqueueTick()
	count := 0
drain:
	for {
		select {
		case <-r.tickCh: count++
		default: break drain
		}
	}
	if count != 1 { t.Fatalf("pending = %d, want 1", count) }
}

func TestRoom_PriorityLane_Full(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	for i := 0; i < PriorityLaneSize; i++ { r.priority <- Envelope{From: "x"} }
	start := time.Now()
	err := r.Leave(context.Background(), "alice")
	elapsed := time.Since(start)
	if !errors.Is(err, ErrPriorityFull) { t.Fatalf("err = %v", err) }
	if elapsed < PriorityLaneWaitTimeout { t.Fatalf("too fast: %v", elapsed) }
}

func TestRoom_PriorityLane_Empty(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	done := make(chan struct{})
	go func() {
		env := <-r.priority
		if env.ReplyTo != nil { env.ReplyTo <- nil }
		close(done)
	}()
	if err := r.Leave(context.Background(), "alice"); err != nil { t.Fatalf("err = %v", err) }
	<-done
}
