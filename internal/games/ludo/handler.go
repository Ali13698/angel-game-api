package ludo

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"time"

	"github.com/angel/angel/internal/core/types"
)

const (
	Seats         = 4
	TokensPerSeat = 4
	Walk          = 52
	Finish        = 57
	HomeColumn    = 6
	MainTrack     = 52
	PosBase       = -1
	MaxSixes      = 3
	extraRoll     = 6
	TurnTimeout   = 30 * time.Second
)

var SafeSquares = map[int]bool{
	0: true, 8: true, 13: true, 21: true,
	26: true, 34: true, 39: true, 47: true,
}

var SeatStart = [Seats]int{0, 13, 26, 39}

type State struct {
	Rev       uint64    `json:"rev"`
	Seats     []string  `json:"seats"`
	Tokens    []int     `json:"tokens"`
	Turn      int       `json:"turn"`
	Die       int       `json:"die"`
	RNG       string    `json:"rng"`
	Moves     int64     `json:"moves"`
	Sixes     int       `json:"sixes"`
	Done      bool      `json:"done"`
	Winner    int       `json:"winner"`
	TurnStart time.Time `json:"turn_start"`
}

func (s *State) Version() uint64 { return s.Rev }
func (s *State) Clone() types.State {
	cp := *s
	cp.Seats = append([]string(nil), s.Seats...)
	cp.Tokens = append([]int(nil), s.Tokens...)
	return &cp
}

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) ID() types.GameID { return "ludo" }

func (h *Handler) NewState(players []types.UserID) types.State {
	s := &State{
		Rev:       1,
		Seats:     make([]string, Seats),
		Tokens:    make([]int, Seats*TokensPerSeat),
		Winner:    -1,
		TurnStart: time.Now(),
	}
	for i := range s.Tokens {
		s.Tokens[i] = PosBase
	}
	for i, p := range players {
		if i < Seats {
			s.Seats[i] = string(p)
		}
	}
	return s
}

func rollDie() int {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return int(binary.LittleEndian.Uint64(b[:])%6) + 1
}

func absolutePos(seat, rel int) int {
	if rel < 0 {
		return -1
	}
	if rel >= MainTrack {
		return -2
	}
	return (SeatStart[seat] + rel) % MainTrack
}

func findKills(s *State, seat, absPos int) []int {
	if absPos < 0 || SafeSquares[absPos] {
		return nil
	}
	var killed []int
	for other := 0; other < Seats; other++ {
		if other == seat || s.Seats[other] == "" {
			continue
		}
		for j := 0; j < TokensPerSeat; j++ {
			idx := other*TokensPerSeat + j
			pos := s.Tokens[idx]
			if pos < 0 || pos >= MainTrack {
				continue
			}
			if absolutePos(other, pos) == absPos {
				killed = append(killed, idx)
			}
		}
	}
	return killed
}

func allFinished(s *State, seat int) bool {
	base := seat * TokensPerSeat
	for i := 0; i < TokensPerSeat; i++ {
		if s.Tokens[base+i] != Finish {
			return false
		}
	}
	return true
}

func nextTurn(s *State) int {
	n := (s.Turn + 1) % Seats
	for i := 0; i < Seats; i++ {
		if s.Seats[n] != "" {
			return n
		}
		n = (n + 1) % Seats
	}
	return n
}

func (h *Handler) ApplyCommand(state types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	s, ok := state.(*State)
	if !ok {
		return state, nil, errors.New("ludo: bad state")
	}
	if s.Done {
		return state, nil, errors.New("ludo: game over")
	}
	seat := -1
	for i, u := range s.Seats {
		if u == string(from) {
			seat = i
			break
		}
	}
	if seat < 0 {
		return state, nil, errors.New("ludo: not seated")
	}
	if seat != s.Turn {
		return state, nil, errors.New("ludo: not your turn")
	}

	next := s.Clone().(*State)
	next.TurnStart = time.Now()

	switch cmd.Type {
	case "roll":
		if next.Die != 0 {
			return state, nil, errors.New("ludo: already rolled")
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
				return next, []types.Event{{
					Type:     "ludo.triple_six",
					Payload:  map[string]any{"seat": seat},
					Audience: types.ToRoom,
				}}, nil
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
		}

		return next, []types.Event{{
			Type:     "ludo.rolled",
			Payload:  map[string]any{"seat": seat, "die": next.Die, "hasMove": hasMove},
			Audience: types.ToRoom,
		}}, nil

	case "move":
		if next.Die == 0 {
			return state, nil, errors.New("ludo: roll first")
		}
		payload, ok := cmd.Payload.(map[string]any)
		if !ok {
			return state, nil, errors.New("ludo: bad payload")
		}
		tokF, ok := payload["token"].(float64)
		if !ok {
			return state, nil, errors.New("ludo: bad token")
		}
		tok := int(tokF)
		if tok < 0 || tok >= TokensPerSeat {
			return state, nil, errors.New("ludo: token out of range")
		}

		idx := seat*TokensPerSeat + tok
		pos := next.Tokens[idx]
		var newPos int
		leaveBase := false

		if pos == PosBase {
			if next.Die != extraRoll {
				return state, nil, errors.New("ludo: need 6 to leave base")
			}
			newPos = 0
			leaveBase = true
		} else {
			newPos = pos + next.Die
			if newPos > Finish {
				return state, nil, errors.New("ludo: overshoot finish")
			}
		}

		next.Tokens[idx] = newPos
		next.Rev++
		next.Moves++

		events := []types.Event{{
			Type: "ludo.moved",
			Payload: map[string]any{
				"seat": seat, "token": tok, "position": newPos, "leaveBase": leaveBase,
			},
			Audience: types.ToRoom,
		}}

		absPos := absolutePos(seat, newPos)
		killed := findKills(next, seat, absPos)
		for _, kidx := range killed {
			next.Tokens[kidx] = PosBase
			events = append(events, types.Event{
				Type:     "ludo.killed",
				Payload:  map[string]any{"tokenIdx": kidx},
				Audience: types.ToRoom,
			})
		}

		if allFinished(next, seat) {
			next.Done = true
			next.Winner = seat
			events = append(events, types.Event{
				Type:     "ludo.win",
				Payload:  map[string]any{"seat": seat},
				Audience: types.ToRoom,
			})
			next.Die = 0
			return next, events, nil
		}

		rolledSix := next.Die == extraRoll
		next.Die = 0
		if !rolledSix {
			next.Sixes = 0
			next.Turn = nextTurn(next)
		}

		return next, events, nil

	default:
		return state, nil, errors.New("ludo: unknown command")
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
	st, ok := s.(*State)
	if !ok {
		return s, nil
	}
	seat := -1
	for i, u := range st.Seats {
		if u == string(id) {
			seat = i
			break
		}
	}
	if seat < 0 {
		return s, nil
	}
	next := st.Clone().(*State)
	next.Seats[seat] = ""
	next.Rev++

	if seat == next.Turn {
		next.Turn = nextTurn(next)
		next.Die = 0
		next.Sixes = 0
	}

	remaining := 0
	lastAlive := -1
	for i, u := range next.Seats {
		if u != "" {
			remaining++
			lastAlive = i
		}
	}
	events := []types.Event{{
		Type:     "ludo.left",
		Payload:  map[string]any{"seat": seat, "remaining": remaining},
		Audience: types.ToRoom,
	}}
	if remaining <= 1 {
		next.Done = true
		if lastAlive >= 0 {
			next.Winner = lastAlive
		}
		events = append(events, types.Event{
			Type:     "ludo.win",
			Payload:  map[string]any{"seat": lastAlive, "reason": "others_left"},
			Audience: types.ToRoom,
		})
	}
	return next, events
}

func (h *Handler) Tick(s types.State, now time.Time) (types.State, []types.Event) {
	st, ok := s.(*State)
	if !ok || st.Done {
		return s, nil
	}
	if st.TurnStart.IsZero() || now.Sub(st.TurnStart) < TurnTimeout {
		return s, nil
	}
	next := st.Clone().(*State)
	next.Turn = nextTurn(next)
	next.Die = 0
	next.Sixes = 0
	next.TurnStart = now
	next.Rev++
	return next, []types.Event{{
		Type:     "ludo.timeout",
		Payload:  map[string]any{"seat": st.Turn},
		Audience: types.ToRoom,
	}}
}

func (h *Handler) IsTerminal(s types.State) bool {
	st, ok := s.(*State)
	return ok && st.Done
}

func (h *Handler) ViewFor(s types.State, viewer types.UserID) any {
	return s
}

var _ types.GameHandler = (*Handler)(nil)
var _ types.StateFilter = (*Handler)(nil)
