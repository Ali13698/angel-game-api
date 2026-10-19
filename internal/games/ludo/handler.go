package ludo

import (
	"errors"
	"time"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
)

func (h *Handler) NewState(players []types.UserID) types.State {
	s := &State{Rev: 1, Seats: make([]string, Seats), Tokens: make([]int, Seats*TokensPerSeat), Winner: -1, TurnStart: time.Now()}
	for i := range s.Tokens { s.Tokens[i] = PosBase }
	for i, p := range players { if i < Seats { s.Seats[i] = string(p) } }
	logger.Info("Ludo state initialized", "players_count", len(players))
	return s
}

func (h *Handler) ApplyCommand(state types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	s, ok := state.(*State)
	if !ok { return state, nil, errors.New("ludo: bad state") }
	if s.Done { return state, nil, errors.New("ludo: game over") }
	
	seat := -1
	for i, u := range s.Seats { if u == string(from) { seat = i; break } }
	if seat < 0 { return state, nil, errors.New("ludo: not seated") }
	if seat != s.Turn { return state, nil, errors.New("ludo: not your turn") }

	next := s.Clone().(*State)
	next.TurnStart = time.Now()

	switch cmd.Type {
	case "roll": return h.handleRoll(next, seat)
	case "move": return h.handleMove(next, seat, cmd)
	default: return state, nil, errors.New("ludo: unknown command")
	}
}

func (h *Handler) OnPlayerJoined(s types.State, id types.UserID) (types.State, []types.Event) {
	st, ok := s.(*State); if !ok { return s, nil }
	for _, u := range st.Seats { if u == string(id) { return s, nil } }
	for i, u := range st.Seats {
		if u == "" {
			n := st.Clone().(*State); n.Seats[i] = string(id); n.Rev++
			return n, []types.Event{{Type:"ludo.joined", Payload:map[string]any{"seat":i,"user_id":string(id)}, Audience:types.ToRoom}}
		}
	}
	return s, nil
}

func (h *Handler) OnPlayerLeft(s types.State, id types.UserID) (types.State, []types.Event) {
	st, ok := s.(*State); if !ok { return s, nil }
	seat := -1
	for i, u := range st.Seats { if u == string(id) { seat = i; break } }
	if seat < 0 { return s, nil }
	
	next := st.Clone().(*State); next.Seats[seat] = ""; next.Rev++
	if seat == next.Turn { next.Turn = nextTurn(next); next.Die = 0; next.Sixes = 0 }
	
	remaining := 0; lastAlive := -1
	for i, u := range next.Seats { if u != "" { remaining++; lastAlive = i } }
	events := []types.Event{{Type: "ludo.left", Payload: map[string]any{"seat": seat, "remaining": remaining}, Audience: types.ToRoom}}
	
	if remaining <= 1 {
		next.Done = true
		if lastAlive >= 0 { next.Winner = lastAlive }
		events = append(events, types.Event{Type: "ludo.win", Payload: map[string]any{"seat": lastAlive, "reason": "others_left"}, Audience: types.ToRoom})
	}
	return next, events
}

func (h *Handler) Tick(s types.State, now time.Time) (types.State, []types.Event) {
	st, ok := s.(*State)
	if !ok || st.Done { return s, nil }
	if st.TurnStart.IsZero() || now.Sub(st.TurnStart) < TurnTimeout { return s, nil }
	
	next := st.Clone().(*State); next.Turn = nextTurn(next); next.Die = 0; next.Sixes = 0; next.TurnStart = now; next.Rev++
	return next, []types.Event{{Type: "ludo.timeout", Payload: map[string]any{"seat": st.Turn}, Audience: types.ToRoom}}
}

func (h *Handler) IsTerminal(s types.State) bool { st, ok := s.(*State); return ok && st.Done }

func (h *Handler) GetResult(s types.State) types.GameResult {
	st, ok := s.(*State)
	if !ok || !st.Done || st.Winner < 0 { return types.GameResult{IsDraw: true} }
	winnerID := types.UserID(st.Seats[st.Winner])
	var loserID types.UserID
	for i, seat := range st.Seats { if i != st.Winner && seat != "" { loserID = types.UserID(seat); break } }
	return types.GameResult{WinnerID: winnerID, LoserID: loserID, IsDraw: false}
}

func (h *Handler) ViewFor(s types.State, viewer types.UserID) any { return s }

var _ types.GameHandler = (*Handler)(nil)
var _ types.StateFilter = (*Handler)(nil)
