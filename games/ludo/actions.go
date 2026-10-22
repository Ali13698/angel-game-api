package main

import "errors"

func handleRoll(next *State, seat int) ([]Event, error) {
	if next.Die != 0 {
		return nil, errors.New("already rolled")
	}
	next.Die = rollDie()
	next.Rev++
	next.Moves++

	if next.Die == extraRoll {
		next.Sixes++
		if next.Sixes >= MaxSixes {
			next.Die = 0
			next.Sixes = 0
			next.Turn = nextTurn(next)
			return []Event{
				{Type: "ludo.triple_six", Payload: map[string]any{"seat": seat}, Audience: 0},
				{Type: "turn_changed", Payload: map[string]any{"seat": next.Turn, "user_id": next.Seats[next.Turn]}, Audience: 0},
			}, nil
		}
	}

	hasMove := false
	base := seat * TokensPerSeat
	for j := 0; j < TokensPerSeat; j++ {
		pos := next.Tokens[base+j]
		if pos == PosBase {
			if next.Die == extraRoll {
				hasMove = true
				break
			}
		} else if pos < Finish && pos+next.Die <= Finish {
			hasMove = true
			break
		}
	}

	if !hasMove && next.Die != extraRoll {
		next.Die = 0
		next.Sixes = 0
		next.Turn = nextTurn(next)
		return []Event{
			{Type: "turn_changed", Payload: map[string]any{"seat": next.Turn, "user_id": next.Seats[next.Turn], "skipped": true}, Audience: 0},
		}, nil
	}

	return []Event{{Type: "ludo.rolled", Payload: map[string]any{"seat": seat, "die": next.Die, "hasMove": hasMove}, Audience: 0}}, nil
}

func handleMove(next *State, seat int, payload map[string]any) ([]Event, error) {
	if next.Die == 0 {
		return nil, errors.New("roll first")
	}
	tokF, ok := payload["token"].(float64)
	if !ok {
		return nil, errors.New("bad token")
	}
	tok := int(tokF)
	if tok < 0 || tok >= TokensPerSeat {
		return nil, errors.New("token out of range")
	}

	idx := seat*TokensPerSeat + tok
	pos := next.Tokens[idx]
	var newPos int
	leaveBase := false

	if pos == PosBase {
		if next.Die != extraRoll {
			return nil, errors.New("need 6 to leave base")
		}
		newPos = 0
		leaveBase = true
	} else {
		newPos = pos + next.Die
		if newPos > Finish {
			return nil, errors.New("overshoot finish")
		}
	}

	next.Tokens[idx] = newPos
	next.Rev++
	next.Moves++
	events := []Event{{Type: "ludo.moved", Payload: map[string]any{"seat": seat, "token": tok, "position": newPos, "leaveBase": leaveBase}, Audience: 0}}

	absPos := absolutePos(seat, newPos)
	killed := findKills(next, seat, absPos)
	for _, kidx := range killed {
		next.Tokens[kidx] = PosBase
		events = append(events, Event{Type: "ludo.killed", Payload: map[string]any{"tokenIdx": kidx}, Audience: 0})
	}

	if allFinished(next, seat) {
		next.Done = true
		next.Winner = seat
		winnerID := next.Seats[seat]
		var loserID string
		for i, s := range next.Seats {
			if i != seat && s != "" {
				loserID = s
				break
			}
		}
		events = append(events, Event{
			Type:     "ludo.win",
			Payload:  map[string]any{"seat": seat, "winner_id": winnerID, "loser_id": loserID},
			Audience: 0,
		})
		next.Die = 0
		return events, nil
	}

	rolledSix := next.Die == extraRoll
	next.Die = 0
	if !rolledSix {
		next.Sixes = 0
		next.Turn = nextTurn(next)
		events = append(events, Event{Type: "turn_changed", Payload: map[string]any{"seat": next.Turn, "user_id": next.Seats[next.Turn]}, Audience: 0})
	}
	return events, nil
}

func ApplyCommand(s State, from string, cmdType string, payload map[string]any) (State, []Event, error) {
	if s.Done {
		return s, nil, errors.New("game over")
	}
	seat := -1
	for i, u := range s.Seats {
		if u == from {
			seat = i
			break
		}
	}
	if seat < 0 {
		return s, nil, errors.New("not seated")
	}
	if seat != s.Turn {
		return s, nil, errors.New("not your turn")
	}

	next := s
	var events []Event
	var err error

	switch cmdType {
	case "roll":
		events, err = handleRoll(&next, seat)
	case "move":
		events, err = handleMove(&next, seat, payload)
	default:
		return s, nil, errors.New("unknown command")
	}

	if err != nil {
		return s, nil, err
	}
	return next, events, nil
}
