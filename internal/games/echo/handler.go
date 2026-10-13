package echo

import (
	"errors"
	"time"

	"github.com/angel/angel/internal/core/types"
)

type State struct {
	Rev  uint64 `json:"rev"`
	Last any    `json:"last,omitempty"`
}

func (s *State) Version() uint64 { return s.Rev }
func (s *State) Clone() types.State {
	cp := *s
	return &cp
}

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) ID() types.GameID { return "echo" }

func (h *Handler) NewState(players []types.UserID) types.State {
	return &State{Rev: 1}
}

func (h *Handler) ApplyCommand(s types.State, from types.UserID, cmd types.Command) (types.State, []types.Event, error) {
	st, ok := s.(*State)
	if !ok {
		return s, nil, errors.New("echo: bad state")
	}
	next := st.Clone().(*State)
	next.Rev++
	next.Last = cmd.Payload
	return next, []types.Event{{
		Type:     "echo.reply",
		Payload:  cmd.Payload,
		Audience: types.ToPlayer,
		Target:   from,
	}}, nil
}

func (h *Handler) OnPlayerJoined(s types.State, id types.UserID) (types.State, []types.Event) { return s, nil }

func (h *Handler) OnPlayerLeft(s types.State, id types.UserID) (types.State, []types.Event) {
	return s, nil
}

func (h *Handler) Tick(s types.State, now time.Time) (types.State, []types.Event) {
	return s, nil
}

// IsTerminal برای echo همیشه false است — این بازی تستی پایانی ندارد.
func (h *Handler) IsTerminal(s types.State) bool { return false }

func (h *Handler) ViewFor(s types.State, viewer types.UserID) any { return s }

var _ types.GameHandler = (*Handler)(nil)
var _ types.StateFilter = (*Handler)(nil)
