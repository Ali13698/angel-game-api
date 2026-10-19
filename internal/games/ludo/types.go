package ludo

import (
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
