package types

import (
	"encoding/json"
	"strconv"
	"time"
)

type UserID string
type RoomID string
type GameID string

// Audience مشخص می‌کند یک Event برای چه کسی ارسال می‌شود.
//
// مقادیر مجاز:
//   ToRoom       — برای همهٔ اعضای اتاق
//   ToPlayer     — برای یک بازیکن مشخص (Target)
//   ToSpectators — برای تماشاگران
//   Internal     — صرفاً داخلی، بدون ارسال به بیرون
type Audience int

const (
	ToRoom Audience = iota
	ToPlayer
	ToSpectators
	Internal
)

type Event struct {
	Type     string   `json:"type"`
	Payload  any      `json:"payload,omitempty"`
	Audience Audience `json:"-"`
	Target   UserID   `json:"-"`
}

func (e Event) Broadcast() bool {
	return e.Audience == ToRoom || e.Audience == ToSpectators
}

type Command struct {
	ID      string `json:"id,omitempty"`
	RoomID  RoomID `json:"room_id,omitempty"`
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
	Version uint64 `json:"version,omitempty"`
}

type GameResult struct {
	WinnerID UserID
	LoserID  UserID
	IsDraw   bool
}

type State interface {
	Version() uint64
	Clone() State
}

// GameHandler رابطی است که هر بازی باید پیاده‌سازی کند.
type GameHandler interface {
	ID() GameID
	NewState(players []UserID) State
	ApplyCommand(s State, from UserID, cmd Command) (State, []Event, error)
	OnPlayerJoined(s State, id UserID) (State, []Event)
	OnPlayerLeft(s State, id UserID) (State, []Event)
	Tick(s State, now time.Time) (State, []Event)
	IsTerminal(s State) bool
	GetResult(s State) GameResult
}

// Game نام قدیمی رابط GameHandler است؛ برای سازگاری موقت.
//
// Deprecated: از GameHandler استفاده کنید.
type Game = GameHandler

type StateFilter interface {
	ViewFor(s State, viewer UserID) any
}

func StrOf(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case UserID:
		return string(x)
	case nil:
		return ""
	default:
		return ""
	}
}

func StrsOf(v any) []string {
	switch x := v.(type) {
	case []string:
		return append([]string(nil), x...)
	case []any:
		out := make([]string, 0, len(x))
		for _, e := range x {
			out = append(out, StrOf(e))
		}
		return out
	}
	return nil
}

func IntOf(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	case json.Number:
		if n, err := x.Int64(); err == nil {
			return int(n)
		}
	case string:
		if n, err := strconv.Atoi(x); err == nil {
			return n
		}
	}
	return 0
}
