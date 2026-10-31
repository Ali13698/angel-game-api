package wire

import (
	"strconv"

	"github.com/angel/angel/internal/core/types"
	"google.golang.org/protobuf/proto"
)

const (
	TypeJoin   = "join"
	TypeLeave  = "leave"
	TypeState  = "state"
	TypePing   = "ping"
	TypeJoined = "joined"
	TypeLeft   = "left"
	TypeEvent  = "event"
	TypePong   = "pong"
	TypeError  = "error"
	TypeChat   = "chat"
)

const (
	CodeBadFrame     = "bad_frame"
	CodeNoRoom       = "no_room"
	CodeRoomFull     = "room_full"
	CodeNotJoined    = "not_joined"
	CodeAlreadyIn    = "already_joined"
	CodeServerClosed = "server_closed"
	CodeSlowClient   = "slow_client"
	CodeInternal     = "internal"
	CodeRejected     = "rejected"
	CodeStaleVersion = "stale_version"
	CodeWrongPhase   = "wrong_phase"
	CodeBusy         = "busy"
	CodeBackpressure = "backpressure"
	CodeCapacityFull = "capacity_full"
)

func EventID(streamID string, seq uint64) string {
	return streamID + ":" + strconv.FormatUint(seq, 10)
}

type OutboundMsg struct {
	Type    string
	Code    string
	RoomID  string
	GameID  string
	Rev     uint64
	Events  []*GameEvent
	State   []byte
	From    string
	Error   string
	Payload []byte
}

func (m OutboundMsg) Bytes() []byte {
	sm := &ServerMsg{
		Type:    m.Type,
		RoomId:  m.RoomID,
		GameId:  m.GameID,
		Rev:     int64(m.Rev),
		From:    m.From,
		State:   m.State,
		Payload: m.Payload,
	}
	if len(m.Events) > 0 {
		sm.Events = m.Events
	}
	if m.Code != "" || m.Error != "" {
		sm.Error = &ErrorMsg{Code: m.Code, Error: m.Error}
	}
	b, err := proto.Marshal(sm)
	if err != nil {
		return nil
	}
	return b
}

func ErrFrame(code, msg string) []byte {
	return OutboundMsg{Type: TypeError, Code: code, Error: msg}.Bytes()
}

type InboundMsg struct {
	ID      string
	Type    string
	RoomID  string
	GameID  string
	Rev     uint64
	Payload any
}

func NewInboundMsg(data []byte) (InboundMsg, error) {
	var cm ClientMsg
	if err := proto.Unmarshal(data, &cm); err != nil {
		return InboundMsg{}, err
	}
	payload, err := DecodeValue(cm.Payload)
	if err != nil {
		return InboundMsg{}, err
	}
	return InboundMsg{
		ID:      cm.Id,
		Type:    cm.Type,
		RoomID:  cm.RoomId,
		GameID:  cm.GameId,
		Rev:     uint64(cm.Rev),
		Payload: payload,
	}, nil
}

func (m InboundMsg) Command() (types.Command, error) {
	cmd := types.Command{ID: m.ID, RoomID: types.RoomID(m.RoomID), Type: m.Type, Version: m.Rev}
	if m.Payload == nil {
		return cmd, nil
	}
	cmd.Payload = m.Payload
	return cmd, nil
}
