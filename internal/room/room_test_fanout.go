package room

import (
	"context"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

func TestRoom_1A_EventOrder_100(t *testing.T) {
	events := make([]types.Event, 100)
	for i := range events { events[i] = types.Event{Type: "e", Payload: i, Audience: types.ToRoom} }
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice", "bob"}, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	activate(t, r)
	_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}})
	calls := waitForCalls(t, tr, 100, 2*time.Second)
	if len(calls) != 100 { t.Fatalf("got %d calls, want 100", len(calls)) }
	for i := 0; i < 100; i++ {
		if len(calls[i].frame.Events) != 1 || calls[i].frame.Events[0] == nil { t.Fatalf("call %d: want one non-nil event", i) }
		if calls[i].frame.Events[0].Seq != int64(i+1) { t.Fatalf("call %d: seq = %d", i, calls[i].frame.Events[0].Seq) }
	}
}

func TestRoom_1A_AudienceRouting(t *testing.T) {
	events := []types.Event{{Type: "to-room", Audience: types.ToRoom}, {Type: "to-player", Audience: types.ToPlayer, Target: "bob"}, {Type: "to-spec", Audience: types.ToSpectators}, {Type: "internal", Audience: types.Internal}}
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice", "bob"}, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	activate(t, r)
	_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}})
	_ = waitForCalls(t, tr, 3, time.Second)
	time.Sleep(100 * time.Millisecond)
	calls := tr.snapshot()
	if len(calls) != 3 { t.Fatalf("got %d calls, want 3", len(calls)) }
	for i, call := range calls {
		if len(call.frame.Events) != 1 || call.frame.Events[0] == nil { t.Fatalf("call %d: want one non-nil event", i) }
	}
	if calls[0].kind != "broadcast" || calls[0].frame.Events[0].Type != "to-room" { t.Fatalf("call 0 failed") }
	if calls[1].kind != "send" || calls[1].target != "bob" || calls[1].frame.Events[0].Type != "to-player" { t.Fatalf("call 1 failed") }
	if calls[2].kind != "broadcast" || calls[2].frame.Events[0].Type != "to-spec" { t.Fatalf("call 2 failed") }
}
