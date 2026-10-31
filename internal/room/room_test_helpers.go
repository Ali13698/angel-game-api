package room

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/wire"
	"google.golang.org/protobuf/proto"
)

type testState struct{ rev uint64 }

func (s *testState) Version() uint64 { return s.rev }
func (s *testState) Clone() types.State { cp := *s; return &cp }

type testGame struct {
	id          types.GameID
	terminal    atomic.Bool
	emitOnApply []types.Event
}

func (g *testGame) ID() types.GameID { return g.id }
func (g *testGame) NewState(_ []types.UserID) types.State { return &testState{rev: 1} }
func (g *testGame) ApplyCommand(s types.State, _ types.UserID, _ types.Command) (types.State, []types.Event, error) {
	st := s.(*testState)
	next := st.Clone().(*testState)
	next.rev++
	return next, g.emitOnApply, nil
}
func (g *testGame) OnPlayerJoined(s types.State, _ types.UserID) (types.State, []types.Event) { return s, nil }
func (g *testGame) OnPlayerLeft(s types.State, _ types.UserID) (types.State, []types.Event) { return s, nil }
func (g *testGame) Tick(s types.State, _ time.Time) (types.State, []types.Event) { return s, nil }
func (g *testGame) IsTerminal(_ types.State) bool { return g.terminal.Load() }
func (g *testGame) GetResult(_ types.State) types.GameResult { return types.GameResult{IsDraw: true} }

type nullTransport struct{}
func (nullTransport) Broadcast(_ []byte) {}
func (nullTransport) SendTo(_ types.UserID, _ []byte) {}
func (nullTransport) CloseSlowClient(_ types.UserID, _ string, _ string) {}
func (nullTransport) SignalBackpressure(_ string, _ string) {}

type recordedCall struct {
	kind   string
	target types.UserID
	code   string
	frame  wire.OutboundMsg
}

type recordingTransport struct {
	mu    sync.Mutex
	calls []recordedCall
}

func (t *recordingTransport) record(c recordedCall) {
	t.mu.Lock()
	t.calls = append(t.calls, c)
	t.mu.Unlock()
}

func decode(frame []byte) (wire.OutboundMsg, bool) {
	var sm wire.ServerMsg
	if err := proto.Unmarshal(frame, &sm); err != nil {
		return wire.OutboundMsg{}, false
	}
	return wire.OutboundMsg{Type: sm.Type, Events: sm.Events}, true
}

func (t *recordingTransport) Broadcast(frame []byte) {
	if msg, ok := decode(frame); ok { t.record(recordedCall{kind: "broadcast", frame: msg}) }
}
func (t *recordingTransport) SendTo(id types.UserID, frame []byte) {
	if msg, ok := decode(frame); ok { t.record(recordedCall{kind: "send", target: id, frame: msg}) }
}
func (t *recordingTransport) CloseSlowClient(id types.UserID, code string, _ string) {
	t.record(recordedCall{kind: "close_slow_client", target: id, code: code})
}
func (t *recordingTransport) SignalBackpressure(code string, _ string) {
	t.record(recordedCall{kind: "backpressure", code: code})
}
func (t *recordingTransport) snapshot() []recordedCall {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]recordedCall(nil), t.calls...)
}
func (t *recordingTransport) countKind(kind string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := 0
	for _, c := range t.calls { if c.kind == kind { n++ } }
	return n
}
func waitForCalls(t *testing.T, tr *recordingTransport, want int, timeout time.Duration) []recordedCall {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		calls := tr.snapshot()
		if len(calls) >= want { return calls }
		time.Sleep(5 * time.Millisecond)
	}
	return tr.snapshot()
}
func activate(t *testing.T, r *Room) {
	t.Helper()
	if err := r.Join(context.Background(), "alice"); err != nil { t.Fatalf("join alice: %v", err) }
	if err := r.Join(context.Background(), "bob"); err != nil { t.Fatalf("join bob: %v", err) }
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && r.Phase() != PhaseActive {
		r.enqueueTick()
		time.Sleep(5 * time.Millisecond)
	}
	if r.Phase() != PhaseActive { t.Fatalf("phase = %v, want PhaseActive", r.Phase()) }
}
func newTestRoom(t *testing.T, inbox int, tr Transport, events []types.Event) *Room {
	t.Helper()
	g := &testGame{id: "test", emitOnApply: events}
	return New(Config{ID: "room-1", GameID: "test", Game: g, Transport: tr, Players: []types.UserID{"alice", "bob"}, InboxSize: inbox, CountdownDuration: 10 * time.Millisecond})
}
func newRoom(g *testGame, tr Transport, players []types.UserID, inbox int) *Room {
	return New(Config{ID: "room-1", GameID: "test", Game: g, Transport: tr, Players: players, InboxSize: inbox, CountdownDuration: 5 * time.Millisecond})
}
