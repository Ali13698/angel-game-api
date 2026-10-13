package ludo

import (
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
)

func TestNewGame(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	if s.Rev != 1 { t.Fatal("bad rev") }
	if s.Seats[0] != "a" || s.Seats[1] != "b" { t.Fatal("bad seats") }
	for _, p := range s.Tokens {
		if p != PosBase { t.Fatal("token not in base") }
	}
}

func TestNeedSixToLeaveBase(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	s.Die = 3
	_, _, err := h.ApplyCommand(s, "a", types.Command{
		Type: "move", Payload: map[string]any{"token": float64(0)},
	})
	if err == nil { t.Fatal("expected error") }
}

func TestRollRandomness(t *testing.T) {
	seen := map[int]bool{}
	for i := 0; i < 500; i++ { seen[rollDie()] = true }
	if len(seen) < 6 { t.Fatalf("rollDie not random: %d values", len(seen)) }
}

func TestTurnTimeout(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	s.TurnStart = time.Now().Add(-TurnTimeout - time.Second)
	next, events := h.Tick(s, time.Now())
	if next.(*State).Turn != 1 { t.Fatal("turn not passed") }
	if len(events) != 1 { t.Fatal("expected 1 event") }
}

func TestPlayerLeft(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	s.Turn = 0
	next, events := h.OnPlayerLeft(s, "a")
	ns := next.(*State)
	if ns.Seats[0] != "" { t.Fatal("seat 0 should be empty") }
	if ns.Turn != 1 { t.Fatal("turn should pass") }
	if !ns.Done { t.Fatal("last player should win") }
	if ns.Winner != 1 { t.Fatal("expected winner 1") }
	if len(events) < 2 { t.Fatal("expected left + win events") }
}

func TestExtraTurnOnSix(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	s.Tokens[0] = 5
	s.Turn = 0
	s.Die = 6
	next, _, err := h.ApplyCommand(s, "a", types.Command{
		Type: "move", Payload: map[string]any{"token": float64(0)},
	})
	if err != nil { t.Fatal(err) }
	ns := next.(*State)
	if ns.Turn != 0 { t.Fatal("should keep turn after 6") }
}

func TestWinCondition(t *testing.T) {
	h := New()
	s := h.NewState([]types.UserID{"a", "b"}).(*State)
	s.Turn = 0
	s.Die = 1
	s.Tokens[0] = Finish
	s.Tokens[1] = Finish
	s.Tokens[2] = Finish
	s.Tokens[3] = Finish - 1
	next, events, err := h.ApplyCommand(s, "a", types.Command{
		Type: "move", Payload: map[string]any{"token": float64(3)},
	})
	if err != nil { t.Fatal(err) }
	ns := next.(*State)
	if !ns.Done { t.Fatal("expected done") }
	if ns.Winner != 0 { t.Fatal("expected winner 0") }
	if !h.IsTerminal(ns) { t.Fatal("IsTerminal should be true") }
	hasWin := false
	for _, e := range events {
		if e.Type == "ludo.win" { hasWin = true }
	}
	if !hasWin { t.Fatal("expected win event") }
}
