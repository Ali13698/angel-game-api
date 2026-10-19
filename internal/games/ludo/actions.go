package ludo

import (
	"errors"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/logger"
)

func (h *Handler) handleRoll(next *State, seat int) (types.State, []types.Event, error) {
	if next.Die != 0 { return next, nil, errors.New("ludo: already rolled") }
	next.Die = rollDie(); next.Rev++; next.Moves++
	logger.Info("Ludo roll", "seat", seat, "die", next.Die)

	if next.Die == extraRoll {
		next.Sixes++
		if next.Sixes >= MaxSixes {
			next.Die = 0; next.Sixes = 0; next.Turn = nextTurn(next)
			return next, []types.Event{{Type: "ludo.triple_six", Payload: map[string]any{"seat": seat}, Audience: types.ToRoom}}, nil
		}
	}

	hasMove := false; base := seat * TokensPerSeat
	for j := 0; j < TokensPerSeat; j++ {
		pos := next.Tokens[base+j]
		if pos == PosBase { if next.Die == extraRoll { hasMove = true; break } } else if pos < Finish && pos+next.Die <= Finish { hasMove = true; break }
	}

	if !hasMove && next.Die != extraRoll { next.Die = 0; next.Sixes = 0; next.Turn = nextTurn(next) }
	return next, []types.Event{{Type: "ludo.rolled", Payload: map[string]any{"seat": seat, "die": next.Die, "hasMove": hasMove}, Audience: types.ToRoom}}, nil
}

func (h *Handler) handleMove(next *State, seat int, cmd types.Command) (types.State, []types.Event, error) {
	if next.Die == 0 { return next, nil, errors.New("ludo: roll first") }
	payload, ok := cmd.Payload.(map[string]any)
	if !ok { return next, nil, errors.New("ludo: bad payload") }
	
	tokF, ok := payload["token"].(float64)
	if !ok { return next, nil, errors.New("ludo: bad token") }
	tok := int(tokF)
	if tok < 0 || tok >= TokensPerSeat { return next, nil, errors.New("ludo: token out of range") }

	idx := seat*TokensPerSeat + tok; pos := next.Tokens[idx]; var newPos int; leaveBase := false

	if pos == PosBase {
		if next.Die != extraRoll { return next, nil, errors.New("ludo: need 6 to leave base") }
		newPos = 0; leaveBase = true
	} else {
		newPos = pos + next.Die
		if newPos > Finish { return next, nil, errors.New("ludo: overshoot finish") }
	}

	next.Tokens[idx] = newPos; next.Rev++; next.Moves++
	events := []types.Event{{Type: "ludo.moved", Payload: map[string]any{"seat": seat, "token": tok, "position": newPos, "leaveBase": leaveBase}, Audience: types.ToRoom}}

	absPos := absolutePos(seat, newPos); killed := findKills(next, seat, absPos)
	for _, kidx := range killed {
		next.Tokens[kidx] = PosBase
		events = append(events, types.Event{Type: "ludo.killed", Payload: map[string]any{"tokenIdx": kidx}, Audience: types.ToRoom})
	}

	if allFinished(next, seat) {
		next.Done = true; next.Winner = seat
		logger.Info("Ludo finished", "winner_seat", seat)
		events = append(events, types.Event{Type: "ludo.win", Payload: map[string]any{"seat": seat}, Audience: types.ToRoom})
		next.Die = 0; return next, events, nil
	}

	rolledSix := next.Die == extraRoll; next.Die = 0
	if !rolledSix { next.Sixes = 0; next.Turn = nextTurn(next) }
	return next, events, nil
}
