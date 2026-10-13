package echo

import (
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

func TestHandler_ID(t *testing.T) {
	h := New()
	if got := h.ID(); got != types.GameID("echo") {
		t.Fatalf("ID() = %q, want %q", got, "echo")
	}
}

func TestHandler_NewState(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"alice", "bob"})
	if s == nil {
		t.Fatal("NewState returned nil")
	}
	if v := s.Version(); v != 1 {
		t.Fatalf("initial Version() = %d, want 1", v)
	}
	if _, ok := s.(*State); !ok {
		t.Fatalf("NewState returned %T, want *State", s)
	}
}

func TestHandler_ApplyCommand(t *testing.T) {
	h := New()
	s := h.NewState(nil)

	from := types.UserID("alice")
	cmd := types.Command{
		ID:      "cmd-1",
		Type:    "echo",
		Payload: "hello",
		Version: 1,
	}

	next, events, err := h.ApplyCommand(s, from, cmd)
	if err != nil {
		t.Fatalf("ApplyCommand error: %v", err)
	}
	if next == nil {
		t.Fatal("next state is nil")
	}
	if next.Version() != s.Version()+1 {
		t.Fatalf("next.Version() = %d, want %d", next.Version(), s.Version()+1)
	}
	if len(events) != 1 {
		t.Fatalf("len(events) = %d, want 1", len(events))
	}

	ev := events[0]
	if ev.Type != "echo.reply" {
		t.Fatalf("event.Type = %q, want %q", ev.Type, "echo.reply")
	}
	if ev.Payload != "hello" {
		t.Fatalf("event.Payload = %v, want %q", ev.Payload, "hello")
	}
	if ev.Audience != types.ToPlayer {
		t.Fatalf("event.Audience = %v, want ToPlayer", ev.Audience)
	}
	if ev.Target != from {
		t.Fatalf("event.Target = %q, want %q", ev.Target, from)
	}
	if ev.Broadcast() {
		t.Fatal("event.Broadcast() = true, want false for ToPlayer")
	}

	// وضعیت اصلی نباید تغییر کند.
	if s.Version() != 1 {
		t.Fatalf("original state mutated: Version() = %d, want 1", s.Version())
	}
	st := s.(*State)
	if st.Last != nil {
		t.Fatalf("original state.Last = %v, want nil", st.Last)
	}
}

func TestHandler_ApplyCommand_BadState(t *testing.T) {
	h := New()
	_, _, err := h.ApplyCommand(nil, "alice", types.Command{})
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestHandler_OnPlayerLeft(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"alice"})

	next, events := h.OnPlayerLeft(s, "alice")
	if next != s {
		t.Fatal("OnPlayerLeft should return the same state")
	}
	if len(events) != 0 {
		t.Fatalf("len(events) = %d, want 0", len(events))
	}
}

func TestHandler_Tick(t *testing.T) {
	h := New()
	s := h.NewState(nil)

	next, events := h.Tick(s, time.Now())
	if next != s {
		t.Fatal("Tick should return the same state")
	}
	if len(events) != 0 {
		t.Fatalf("len(events) = %d, want 0", len(events))
	}
}

func TestHandler_IsTerminal(t *testing.T) {
	h := New()
	s := h.NewState(nil)
	if h.IsTerminal(s) {
		t.Fatal("echo should never be terminal")
	}
	if h.IsTerminal(nil) {
		t.Fatal("IsTerminal(nil) should be false")
	}
}

func TestState_Clone(t *testing.T) {
	orig := &State{Rev: 5, Last: "x"}
	cp := orig.Clone().(*State)

	if cp.Rev != orig.Rev {
		t.Fatalf("Clone().Rev = %d, want %d", cp.Rev, orig.Rev)
	}
	if cp.Last != orig.Last {
		t.Fatalf("Clone().Last = %v, want %v", cp.Last, orig.Last)
	}

	// تغییر در کلون نباید به اصلی اثر بگذارد.
	cp.Rev = 99
	if orig.Rev == 99 {
		t.Fatal("mutating clone affected original")
	}
}

func TestHandler_ViewFor(t *testing.T) {
	h := New()
	s := h.NewState(nil)
	if h.ViewFor(s, "alice") != s {
		t.Fatal("ViewFor should return the state itself")
	}
}

// اطمینان از اینکه Handler هر دو رابط را پیاده‌سازی می‌کند.
func TestHandler_ImplementsInterfaces(t *testing.T) {
	var _ types.GameHandler = (*Handler)(nil)
	var _ types.StateFilter = (*Handler)(nil)
}
