package room

import (
	"context"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/wire"
)

func TestRoom_1A_Overflow_ToPlayer_ClosesClient(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{{Type: "e", Audience: types.ToPlayer, Target: "bob"}}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice", "bob"}, 8)
	for i := 0; i < FanoutBufferSize; i++ { r.fanoutCh <- []pendingEvent{} }
	r.emit(events)
	if r.SlowClientClosed() != 1 { t.Fatalf("SlowClientClosed = %d", r.SlowClientClosed()) }
	if tr.countKind("close_slow_client") != 1 { t.Fatalf("close_slow_client calls = %d", tr.countKind("close_slow_client")) }
	calls := tr.snapshot()
	if len(calls) != 1 || calls[0].target != "bob" || calls[0].code != wire.CodeSlowClient { t.Fatalf("close details wrong") }
	if r.Dropped() != 0 { t.Fatalf("Dropped = %d", r.Dropped()) }
}

func TestRoom_1A_Overflow_Internal_Dropped(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{{Type: "e", Audience: types.Internal}, {Type: "e2", Audience: types.ToSpectators}}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice"}, 8)
	for i := 0; i < FanoutBufferSize; i++ { r.fanoutCh <- []pendingEvent{} }
	r.emit(events)
	if r.Dropped() != 2 { t.Fatalf("Dropped = %d", r.Dropped()) }
	if r.SlowClientClosed() != 0 { t.Fatalf("SlowClientClosed = %d", r.SlowClientClosed()) }
	if r.Backpressure() != 0 { t.Fatalf("Backpressure = %d", r.Backpressure()) }
}

func TestRoom_1A_Overflow_ToRoom_Backpressure(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{{Type: "e", Audience: types.ToRoom}}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice"}, 8)
	for i := 0; i < FanoutBufferSize; i++ { r.fanoutCh <- []pendingEvent{} }
	r.emit(events)
	if r.Backpressure() != 1 { t.Fatalf("Backpressure = %d", r.Backpressure()) }
	if r.Dropped() != 0 { t.Fatalf("Dropped = %d", r.Dropped()) }
	if tr.countKind("backpressure") != 1 { t.Fatalf("backpressure calls = %d", tr.countKind("backpressure")) }
}

func TestRoom_Shutdown_DrainsFanout(t *testing.T) {
	events := make([]types.Event, 100)
	for i := range events { events[i] = types.Event{Type: "e", Audience: types.ToRoom} }
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := newRoom(g, tr, []types.UserID{"alice", "bob"}, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	activate(t, r)
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}}); err != nil { t.Fatalf("Submit err = %v", err) }
	fanoutDeadline := time.Now().Add(2 * time.Second)
	for len(tr.snapshot()) < 100 {
		if time.Now().After(fanoutDeadline) { t.Fatalf("fanout never reached 100 calls") }
		time.Sleep(2 * time.Millisecond)
	}
	if err := r.Close(context.Background()); err != nil { t.Fatalf("Close err = %v", err) }
	if len(tr.snapshot()) != 100 { t.Fatalf("events lost on shutdown") }
}
