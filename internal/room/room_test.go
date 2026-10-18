package room

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/wire"
	"google.golang.org/protobuf/proto"
)

// --- Test doubles ---

type testState struct{ rev uint64 }

func (s *testState) Version() uint64    { return s.rev }
func (s *testState) Clone() types.State { cp := *s; return &cp }

type testGame struct {
	id          types.GameID
	terminal    atomic.Bool
	emitOnApply []types.Event
}

func (g *testGame) ID() types.GameID                       { return g.id }
func (g *testGame) NewState(_ []types.UserID) types.State  { return &testState{rev: 1} }
func (g *testGame) ApplyCommand(s types.State, _ types.UserID, _ types.Command) (types.State, []types.Event, error) {
	st := s.(*testState)
	next := st.Clone().(*testState)
	next.rev++
	return next, g.emitOnApply, nil
}
func (g *testGame) OnPlayerLeft(s types.State, _ types.UserID) (types.State, []types.Event) {
	return s, nil
}
func (g *testGame) Tick(s types.State, _ time.Time) (types.State, []types.Event) {
	return s, nil
}
func (g *testGame) IsTerminal(_ types.State) bool { return g.terminal.Load() }
func (g *testGame) GetResult(_ types.State) types.GameResult { return types.GameResult{IsDraw: true} }

type nullTransport struct{}

func (nullTransport) Broadcast(_ []byte)                                       {}
func (nullTransport) SendTo(_ types.UserID, _ []byte)                          {}
func (nullTransport) CloseSlowClient(_ types.UserID, _ string, _ string)       {}
func (nullTransport) SignalBackpressure(_ string, _ string)                    {}

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

// Broadcast حالا فریم‌های باینری Protobuf را پردازش می‌کند
func (t *recordingTransport) Broadcast(frame []byte) {
	var sm wire.ServerMsg
	if err := proto.Unmarshal(frame, &sm); err == nil {
		var ev *wire.Event
		if sm.Event != nil {
			ev = &wire.Event{
				EventID:  sm.Event.EventId,
				StreamID: sm.Event.StreamId,
				Seq:      uint64(sm.Event.Seq),
				Type:     sm.Event.Type,
			}
		}
		msg := wire.OutboundMsg{Type: sm.Type, Event: ev}
		t.record(recordedCall{kind: "broadcast", frame: msg})
	}
}

func (t *recordingTransport) SendTo(id types.UserID, frame []byte) {
	var sm wire.ServerMsg
	if err := proto.Unmarshal(frame, &sm); err == nil {
		var ev *wire.Event
		if sm.Event != nil {
			ev = &wire.Event{
				EventID:  sm.Event.EventId,
				StreamID: sm.Event.StreamId,
				Seq:      uint64(sm.Event.Seq),
				Type:     sm.Event.Type,
			}
		}
		msg := wire.OutboundMsg{Type: sm.Type, Event: ev}
		t.record(recordedCall{kind: "send", target: id, frame: msg})
	}
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
	for _, c := range t.calls {
		if c.kind == kind {
			n++
		}
	}
	return n
}

func waitForCalls(t *testing.T, tr *recordingTransport, want int, timeout time.Duration) []recordedCall {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		calls := tr.snapshot()
		if len(calls) >= want {
			return calls
		}
		time.Sleep(5 * time.Millisecond)
	}
	return tr.snapshot()
}

func newTestRoom(t *testing.T, inbox int, tr Transport, events []types.Event) *Room {
	t.Helper()
	g := &testGame{id: "test", emitOnApply: events}
	return New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice", "bob"},
		InboxSize:         inbox,
		CountdownDuration: 10 * time.Millisecond,
	})
}

// --- FSM ---

func TestRoom_PhaseTransitions_Allowed(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	if r.Phase() != PhaseLobby {
		t.Fatalf("initial = %v", r.Phase())
	}
	now := time.Now()
	r.advancePhase(PhaseCountdown, now)
	r.advancePhase(PhaseActive, now)
	r.advancePhase(PhaseFinished, now)
	if r.Phase() != PhaseFinished {
		t.Fatalf("final = %v", r.Phase())
	}
}

func TestRoom_PhaseTransitions_Disallowed(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	now := time.Now()
	r.advancePhase(PhaseActive, now)
	if r.Phase() != PhaseLobby {
		t.Fatalf("skip COUNTDOWN should reject, got %v", r.Phase())
	}
	r.advancePhase(PhaseAborted, now)
	r.advancePhase(PhaseCountdown, now)
	if r.Phase() != PhaseAborted {
		t.Fatalf("ABORTED should be terminal, got %v", r.Phase())
	}
}

func TestRoom_CheckTerminal(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	now := time.Now()
	r.advancePhase(PhaseCountdown, now)
	r.advancePhase(PhaseActive, now)
	r.cfg.Game.(*testGame).terminal.Store(true)
	r.checkTerminal()
	if r.Phase() != PhaseFinished {
		t.Fatalf("phase = %v", r.Phase())
	}
}

// --- Inbox overflow ---

func TestRoom_Inbox_Overflow(t *testing.T) {
	r := newTestRoom(t, 2, nullTransport{}, nil)
	for i := 0; i < 2; i++ {
		r.inbox <- Envelope{From: "x"}
	}
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "noop"}}); !errors.Is(err, ErrBusy) {
		t.Fatalf("err = %v, want ErrBusy", err)
	}
}

// --- Tick coalesce ---

func TestRoom_TickSignal_Coalesces(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	r.enqueueTick()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			r.enqueueTick()
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("enqueueTick blocked")
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
		case <-r.tickCh:
			count++
		default:
			break drain
		}
	}
	if count != 1 {
		t.Fatalf("pending = %d, want 1", count)
	}
}

// --- Priority lane ---

func TestRoom_PriorityLane_Full(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	for i := 0; i < PriorityLaneSize; i++ {
		r.priority <- Envelope{From: "x"}
	}
	start := time.Now()
	err := r.Leave(context.Background(), "alice")
	elapsed := time.Since(start)
	if !errors.Is(err, ErrPriorityFull) {
		t.Fatalf("err = %v", err)
	}
	if elapsed < PriorityLaneWaitTimeout {
		t.Fatalf("too fast: %v", elapsed)
	}
}

func TestRoom_PriorityLane_Empty(t *testing.T) {
	r := newTestRoom(t, 8, nullTransport{}, nil)
	done := make(chan struct{})
	go func() {
		env := <-r.priority
		if env.ReplyTo != nil {
			env.ReplyTo <- nil
		}
		close(done)
	}()
	if err := r.Leave(context.Background(), "alice"); err != nil {
		t.Fatalf("err = %v", err)
	}
	<-done
}

// --- 1A: Fanout ---

type slowTransport struct {
	delay time.Duration
	calls atomic.Int64
}

func (s *slowTransport) Broadcast(_ []byte)                                 { time.Sleep(s.delay); s.calls.Add(1) }
func (s *slowTransport) SendTo(_ types.UserID, _ []byte)                    { time.Sleep(s.delay); s.calls.Add(1) }
func (s *slowTransport) CloseSlowClient(_ types.UserID, _ string, _ string) {}
func (s *slowTransport) SignalBackpressure(_ string, _ string)              {}

func TestRoom_1A_HotPath_NotBlockedByFanout(t *testing.T) {
	g := &testGame{
		id: "test",
		emitOnApply: []types.Event{
			{Type: "e1", Audience: types.ToRoom},
			{Type: "e2", Audience: types.ToRoom},
			{Type: "e3", Audience: types.ToRoom},
		},
	}
	slow := &slowTransport{delay: 100 * time.Millisecond}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: slow,
		Players:           []types.UserID{"alice"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && r.Phase() != PhaseActive {
		time.Sleep(5 * time.Millisecond)
	}

	start := time.Now()
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}}); err != nil {
		t.Fatalf("Submit err = %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Fatalf("Submit took %v; hot path blocked by fanout", elapsed)
	}

	deadline = time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && slow.calls.Load() < 3 {
		time.Sleep(10 * time.Millisecond)
	}
	if slow.calls.Load() < 3 {
		t.Fatalf("only %d calls delivered", slow.calls.Load())
	}
}

type multiClientTransport struct {
	mu     sync.Mutex
	delays map[types.UserID]time.Duration
	closed map[types.UserID]bool
	calls  atomic.Int64
}

func (t *multiClientTransport) Broadcast(_ []byte) {
	t.mu.Lock()
	ids := make([]types.UserID, 0, len(t.delays))
	for id := range t.delays {
		ids = append(ids, id)
	}
	t.mu.Unlock()
	for _, id := range ids {
		t.mu.Lock()
		d := t.delays[id]
		t.mu.Unlock()
		time.Sleep(d)
	}
	t.calls.Add(1)
}
func (t *multiClientTransport) SendTo(id types.UserID, _ []byte) {
	t.mu.Lock()
	d := t.delays[id]
	t.mu.Unlock()
	time.Sleep(d)
	t.calls.Add(1)
}
func (t *multiClientTransport) CloseSlowClient(id types.UserID, _ string, _ string) {
	t.mu.Lock()
	t.closed[id] = true
	t.mu.Unlock()
}
func (t *multiClientTransport) SignalBackpressure(_ string, _ string) {}

func TestRoom_1A_HotPath_MultiClient(t *testing.T) {
	g := &testGame{
		id:          "test",
		emitOnApply: []types.Event{{Type: "e", Audience: types.ToRoom}},
	}
	tr := &multiClientTransport{
		delays: map[types.UserID]time.Duration{
			"alice": 10 * time.Millisecond,
			"bob":   10 * time.Millisecond,
			"carol": 500 * time.Millisecond,
		},
		closed: make(map[types.UserID]bool),
	}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice", "bob", "carol"},
		InboxSize:         64,
		CountdownDuration: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && r.Phase() != PhaseActive {
		time.Sleep(5 * time.Millisecond)
	}

	start := time.Now()
	for i := 0; i < 10; i++ {
		_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}})
	}
	elapsed := time.Since(start)
	if elapsed > 50*time.Millisecond {
		t.Fatalf("10 Submits took %v; hot path blocked by slow client", elapsed)
	}
}

func TestRoom_1A_EventOrder_100(t *testing.T) {
	events := make([]types.Event, 100)
	for i := range events {
		events[i] = types.Event{Type: "e", Payload: i, Audience: types.ToRoom}
	}
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && r.Phase() != PhaseActive {
		time.Sleep(5 * time.Millisecond)
	}

	_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}})
	calls := waitForCalls(t, tr, 100, 2*time.Second)
	if len(calls) != 100 {
		t.Fatalf("got %d calls, want 100", len(calls))
	}
	for i := 0; i < 100; i++ {
		if calls[i].frame.Event == nil {
			t.Fatalf("call %d: nil event", i)
		}
		wantSeq := uint64(i + 1)
		if calls[i].frame.Event.Seq != wantSeq {
			t.Fatalf("call %d: seq = %d, want %d", i, calls[i].frame.Event.Seq, wantSeq)
		}
	}
}

func TestRoom_1A_AudienceRouting(t *testing.T) {
	events := []types.Event{
		{Type: "to-room", Audience: types.ToRoom},
		{Type: "to-player", Audience: types.ToPlayer, Target: "bob"},
		{Type: "to-spec", Audience: types.ToSpectators},
		{Type: "internal", Audience: types.Internal},
	}
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice", "bob"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && r.Phase() != PhaseActive {
		time.Sleep(5 * time.Millisecond)
	}

	_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}})
	_ = waitForCalls(t, tr, 3, time.Second)
	time.Sleep(100 * time.Millisecond)
	calls := tr.snapshot()

	if len(calls) != 3 {
		t.Fatalf("got %d calls, want 3 (Internal dropped)", len(calls))
	}
	if calls[0].kind != "broadcast" || calls[0].frame.Event.Type != "to-room" {
		t.Fatalf("call 0: kind=%q type=%q", calls[0].kind, calls[0].frame.Event.Type)
	}
	if calls[1].kind != "send" || calls[1].target != "bob" || calls[1].frame.Event.Type != "to-player" {
		t.Fatalf("call 1: kind=%q target=%q type=%q", calls[1].kind, calls[1].target, calls[1].frame.Event.Type)
	}
	if calls[2].kind != "broadcast" || calls[2].frame.Event.Type != "to-spec" {
		t.Fatalf("call 2: kind=%q type=%q", calls[2].kind, calls[2].frame.Event.Type)
	}
	for i, c := range calls {
		if c.frame.Event.Type == "internal" {
			t.Fatalf("call %d: Internal event leaked", i)
		}
	}
}

// --- 1A: Overflow policy ---

func TestRoom_1A_Overflow_ToPlayer_ClosesClient(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{
		{Type: "e", Audience: types.ToPlayer, Target: "bob"},
	}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice", "bob"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})

	for i := 0; i < FanoutBufferSize; i++ {
		r.fanoutCh <- []pendingEvent{}
	}
	r.emit(events)

	if got := r.SlowClientClosed(); got != 1 {
		t.Fatalf("SlowClientClosed = %d, want 1", got)
	}
	if got := tr.countKind("close_slow_client"); got != 1 {
		t.Fatalf("close_slow_client calls = %d, want 1", got)
	}
	calls := tr.snapshot()
	if len(calls) != 1 {
		t.Fatalf("total calls = %d, want 1", len(calls))
	}
	if calls[0].target != "bob" {
		t.Fatalf("close target = %q, want bob", calls[0].target)
	}
	if calls[0].code != wire.CodeSlowClient {
		t.Fatalf("close code = %q, want %q", calls[0].code, wire.CodeSlowClient)
	}
	if got := r.Dropped(); got != 0 {
		t.Fatalf("Dropped = %d, want 0 (ToPlayer باید رد نشود)", got)
	}
}

func TestRoom_1A_Overflow_Internal_Dropped(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{
		{Type: "e", Audience: types.Internal},
		{Type: "e2", Audience: types.ToSpectators},
	}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})
	for i := 0; i < FanoutBufferSize; i++ {
		r.fanoutCh <- []pendingEvent{}
	}

	r.emit(events)

	if got := r.Dropped(); got != 2 {
		t.Fatalf("Dropped = %d, want 2", got)
	}
	if got := r.SlowClientClosed(); got != 0 {
		t.Fatalf("SlowClientClosed = %d, want 0", got)
	}
	if got := r.Backpressure(); got != 0 {
		t.Fatalf("Backpressure = %d, want 0", got)
	}
}

func TestRoom_1A_Overflow_ToRoom_Backpressure(t *testing.T) {
	tr := &recordingTransport{}
	events := []types.Event{
		{Type: "e", Audience: types.ToRoom},
	}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})
	for i := 0; i < FanoutBufferSize; i++ {
		r.fanoutCh <- []pendingEvent{}
	}

	r.emit(events)

	if got := r.Backpressure(); got != 1 {
		t.Fatalf("Backpressure = %d, want 1", got)
	}
	if got := r.Dropped(); got != 0 {
		t.Fatalf("Dropped = %d, want 0", got)
	}
	if got := tr.countKind("backpressure"); got != 1 {
		t.Fatalf("backpressure calls = %d, want 1", got)
	}
}

// --- بند ۳.۴: shutdown drain ---

func TestRoom_Shutdown_DrainsFanout(t *testing.T) {
	events := make([]types.Event, 100)
	for i := range events {
		events[i] = types.Event{Type: "e", Audience: types.ToRoom}
	}
	tr := &recordingTransport{}
	g := &testGame{id: "test", emitOnApply: events}
	r := New(Config{
		ID: "room-1", GameID: "test", Game: g, Transport: tr,
		Players:           []types.UserID{"alice"},
		InboxSize:         8,
		CountdownDuration: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// فاز فقط داخل tick() منتقل میشود (room.go:421-427)، و منبع تیکِ پیش‌فرض
	// تیکر ۱ ثانیهای است (room.go:302) که با مهلت ۱ ثانیهای تست مسابقه دارد.
	// پس به‌جای انتظارِ زمانی، خودمان یک تیک قطعی تزریق میکنیم (room.go:268-273).
	time.Sleep(20 * time.Millisecond) // > CountdownDuration = 5ms
	r.enqueueTick()
	phaseDeadline := time.Now().Add(2 * time.Second)
	for r.Phase() != PhaseActive {
		if time.Now().After(phaseDeadline) {
			t.Fatalf("phase stuck at %v, want ACTIVE", r.Phase())
		}
		time.Sleep(2 * time.Millisecond)
	}

	// Submit بدون کانال پاسخ است (room.go:177-179)؛ پس نقطهٔ قطعیِ «پردازش شد»
	// مشاهدهٔ رسیدن ۱۰۰ فراخوانی Transport است.
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}}); err != nil {
		t.Fatalf("Submit err = %v", err)
	}
	fanoutDeadline := time.Now().Add(2 * time.Second)
	for len(tr.snapshot()) < 100 {
		if time.Now().After(fanoutDeadline) {
			t.Fatalf("fanout never reached 100 calls (got %d)", len(tr.snapshot()))
		}
		time.Sleep(2 * time.Millisecond)
	}

	if err := r.Close(context.Background()); err != nil {
		t.Fatalf("Close err = %v", err)
	}
	if got := len(tr.snapshot()); got != 100 {
		t.Fatalf("events lost on shutdown: got %d, want 100", got)
	}
}

// --- Concurrency smoke ---

func TestRoom_ConcurrentSubmitAndClose(t *testing.T) {
	r := newTestRoom(t, 64, nullTransport{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "noop"}})
		}()
	}
	wg.Wait()
	if err := r.Close(context.Background()); err != nil {
		t.Fatalf("Close err = %v", err)
	}
}
