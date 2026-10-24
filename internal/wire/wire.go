package wire

import (
	"encoding/json"
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
)

type Event struct {
	EventID  string `json:"event_id"`
	StreamID string `json:"stream_id"`
	Seq      uint64 `json:"seq"`
	Type     string `json:"type"`
	Payload  any    `json:"payload,omitempty"`
}

func EventID(streamID string, seq uint64) string {
	return streamID + ":" + strconv.FormatUint(seq, 10)
}

func NewEvent(streamID string, seq uint64, e types.Event) Event {
	return Event{EventID: EventID(streamID, seq), StreamID: streamID, Seq: seq, Type: e.Type, Payload: e.Payload}
}

type OutboundMsg struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	RoomID  string `json:"room_id,omitempty"`
	GameID  string `json:"game_id,omitempty"`
	Rev     uint64 `json:"rev,omitempty"`
	Events  []*GameEvent `json:"events,omitempty"` // AL V4 fix
	State   any    `json:"state,omitempty"`
	From    string `json:"from,omitempty"`
	Error   string `json:"error,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

func (m OutboundMsg) Bytes() []byte {
	sm := &ServerMsg{Type: m.Type, RoomId: m.RoomID, GameId: m.GameID, Rev: int64(m.Rev), From: m.From}
	
	if len(m.Events) > 0 {
		sm.Events = m.Events
	}
	if m.State != nil {
		s, _ := json.Marshal(m.State)
		sm.State = s
	}
	if m.Payload != nil {
		p, _ := json.Marshal(m.Payload)
		sm.Payload = p
	}
	if m.Code != "" || m.Error != "" {
		sm.Error = &ErrorMsg{Code: m.Code, Error: m.Error}
	}
	b, err := proto.Marshal(sm)
	if err != nil { return nil }
	return b
}

func Event(e *GameEvent) []*GameEvent { return []*GameEvent{e} }

func ErrFrame(code, msg string) []byte {
	return OutboundMsg{Type: TypeError, Code: code, Error: msg}.Bytes()
}

type InboundMsg struct {
	ID      string
	Type    string
	RoomID  string
	GameID  string
	Rev     uint64
	Payload json.RawMessage
}

func NewInboundMsg(data []byte) (InboundMsg, error) {
	var cm ClientMsg
	if err := proto.Unmarshal(data, &cm); err != nil {
		return InboundMsg{}, err
	}
	var p json.RawMessage
	if len(cm.Payload) > 0 {
		p = json.RawMessage(cm.Payload)
	}
	return InboundMsg{ID: cm.Id, Type: cm.Type, RoomID: cm.RoomId, Rev: uint64(cm.Rev), Payload: p}, nil
}

func (m InboundMsg) Command() (types.Command, error) {
	cmd := types.Command{ID: m.ID, RoomID: types.RoomID(m.RoomID), Type: m.Type, Version: m.Rev}
	if len(m.Payload) == 0 || string(m.Payload) == "null" { return cmd, nil }
	var p any
	if err := json.Unmarshal(m.Payload, &p); err != nil { return cmd, err }
	cmd.Payload = p
	return cmd, nil
}
