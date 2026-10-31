package room

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

type slowTransport struct {
	delay time.Duration
	calls atomic.Int64
}
func (s *slowTransport) Broadcast(_ []byte) { time.Sleep(s.delay); s.calls.Add(1) }
func (s *slowTransport) SendTo(_ types.UserID, _ []byte) { time.Sleep(s.delay); s.calls.Add(1) }
func (s *slowTransport) CloseSlowClient(_ types.UserID, _ string, _ string) {}
func (s *slowTransport) SignalBackpressure(_ string, _ string) {}

func TestRoom_1A_HotPath_NotBlockedByFanout(t *testing.T) {
	g := &testGame{id: "test", emitOnApply: []types.Event{{Type: "e1", Audience: types.ToRoom}, {Type: "e2", Audience: types.ToRoom}, {Type: "e3", Audience: types.ToRoom}}}
	slow := &slowTransport{delay: 100 * time.Millisecond}
	r := newRoom(g, slow, []types.UserID{"alice", "bob"}, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	activate(t, r)
	start := time.Now()
	if err := r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}}); err != nil { t.Fatalf("Submit err = %v", err) }
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond { t.Fatalf("Submit took %v", elapsed) }
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && slow.calls.Load() < 3 { time.Sleep(10 * time.Millisecond) }
	if slow.calls.Load() < 3 { t.Fatalf("only %d calls delivered", slow.calls.Load()) }
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
	for id := range t.delays { ids = append(ids, id) }
	t.mu.Unlock()
	for _, id := range ids {
		t.mu.Lock(); d := t.delays[id]; t.mu.Unlock()
		time.Sleep(d)
	}
	t.calls.Add(1)
}
func (t *multiClientTransport) SendTo(id types.UserID, _ []byte) {
	t.mu.Lock(); d := t.delays[id]; t.mu.Unlock()
	time.Sleep(d); t.calls.Add(1)
}
func (t *multiClientTransport) CloseSlowClient(id types.UserID, _ string, _ string) {
	t.mu.Lock(); t.closed[id] = true; t.mu.Unlock()
}
func (t *multiClientTransport) SignalBackpressure(_ string, _ string) {}

func TestRoom_1A_HotPath_MultiClient(t *testing.T) {
	g := &testGame{id: "test", emitOnApply: []types.Event{{Type: "e", Audience: types.ToRoom}}}
	tr := &multiClientTransport{delays: map[types.UserID]time.Duration{"alice": 10 * time.Millisecond, "bob": 10 * time.Millisecond, "carol": 500 * time.Millisecond}, closed: make(map[types.UserID]bool)}
	r := newRoom(g, tr, []types.UserID{"alice", "bob", "carol"}, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	activate(t, r)
	start := time.Now()
	for i := 0; i < 10; i++ { _ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "x"}}) }
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond { t.Fatalf("10 Submits took %v", elapsed) }
}

func TestRoom_ConcurrentSubmitAndClose(t *testing.T) {
	r := newTestRoom(t, 64, nullTransport{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = r.Submit(Request{From: "alice", Cmd: types.Command{Type: "noop"}}) }()
	}
	wg.Wait()
	if err := r.Close(context.Background()); err != nil { t.Fatalf("Close err = %v", err) }
}
