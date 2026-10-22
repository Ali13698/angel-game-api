package main

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
)

var SafeSquares = map[int]bool{
	0: true, 8: true, 13: true, 21: true,
	26: true, 34: true, 39: true, 47: true,
}

var SeatStart = [Seats]int{0, 13, 26, 39}

type State struct {
	Rev    uint64   `json:"rev"`
	Seats  []string `json:"seats"`
	Tokens []int    `json:"tokens"`
	Turn   int      `json:"turn"`
	Die    int      `json:"die"`
	RNG    string   `json:"rng"`
	Moves  int64    `json:"moves"`
	Sixes  int      `json:"sixes"`
	Done   bool     `json:"done"`
	Winner int      `json:"winner"`
}

type Event struct {
	Type     string `json:"type"`
	Payload  any    `json:"payload"`
	Audience int    `json:"audience"` // 0=ToRoom
	Target   string `json:"target,omitempty"`
}

type Intent struct {
	From string `json:"from"`
	Cmd  struct {
		Type    string         `json:"type"`
		Payload map[string]any `json:"payload"`
	} `json:"cmd"`
}
