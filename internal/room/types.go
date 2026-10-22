package room

import (
	"context"
	"errors"
	"time"

	"github.com/angel/angel/internal/core/types"
)

var (
	ErrClosed        = errors.New("room: closed")
	ErrBusy          = errors.New("room: inbox is full")
	ErrWrongPhase    = errors.New("room: command not accepted in this phase")
	ErrStaleVersion  = errors.New("room: stale command version")
	ErrPriorityFull  = errors.New("room: priority lane full")
)

type Phase int

const (
	PhaseLobby Phase = iota
	PhaseCountdown
	PhaseActive
	PhaseFinished
	PhaseAborted
)

func (p Phase) String() string {
	switch p {
	case PhaseLobby:
		return "LOBBY"
	case PhaseCountdown:
		return "COUNTDOWN"
	case PhaseActive:
		return "ACTIVE"
	case PhaseFinished:
		return "FINISHED"
	case PhaseAborted:
		return "ABORTED"
	}
	return "UNKNOWN"
}

const (
	DefaultCountdownDuration = 3 * time.Second
	PriorityLaneSize         = 64
	PriorityLaneWaitTimeout  = 100 * time.Millisecond
	TickBufferSize           = 1
	FanoutBufferSize         = 256
)

type Transport interface {
	Broadcast([]byte)
	SendTo(types.UserID, []byte)
	CloseSlowClient(types.UserID, string, string)
	SignalBackpressure(string, string)
}

type StatsRecorder interface {
	RecordMatch(context.Context, types.UserID, string, bool, bool) error
}

// EventPublisher رابط برای ارسال رویدادها به کافکا
type EventPublisher interface {
	Publish(ctx context.Context, topic string, key string, payload any) error
}

type Config struct {
	ID                types.RoomID
	GameID            types.GameID
	Game              types.Game
	Transport         Transport
	Players           []types.UserID
	InboxSize         int
	CountdownDuration time.Duration
	StatsRecorder     StatsRecorder
	EventPublisher    EventPublisher // تزریق کافکا
}

type Request struct {
	From types.UserID
	Cmd  types.Command
}

type Envelope struct {
	From    types.UserID
	Cmd     types.Command
	Leave   bool
	ReplyTo chan<- error
}

type pendingEvent struct {
	seq      uint64
	event    types.Event
	streamID string
	gameID   string
	rev      uint64
}
