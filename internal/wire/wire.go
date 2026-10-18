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
	TypeEmoji  = "emoji"
)

const (
	CodeBadFrame     = "bad_frame"
	CodeNoRoom       = "no_room"
	CodeRoomFull     = "room_full"
	CodeNotJoined    = "not_joined"
	CodeAlreadyIn    = "already_joined"
	CodeUnknownGame  = "unknown_game"
	CodeServerClosed = "server_closed"
	CodeSlowConsumer = "slow_consumer"
	CodeSlowClient   = "slow_client"
	CodeInternal     = "internal"
	CodeRejected     = "rejected"
	CodeTimeout      = "timeout"
	CodeBusy         = "busy"
	CodeStaleVersion = "stale_version"
	CodeWrongPhase   = "wrong_phase"
	CodeBackpressure = "backpressure"
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
	return Event{
		EventID:  EventID(streamID, seq),
		StreamID: streamID,
		Seq:      seq,
		Type:     e.Type,
		Payload:  e.Payload,
	}
}

type OutboundMsg struct {
	Type   string `json:"type"`
	Code   string `json:"code,omitempty"`
	RoomID string `json:"room_id,omitempty"`
	GameID string `json:"game_id,omitempty"`
	Rev    uint64 `json:"rev,omitempty"`
	Event  *Event `json:"event,omitempty"`
	State  any    `json:"state,omitempty"`
	From   string `json:"from,omitempty"`
	Error  string `json:"error,omitempty"`
	Payload any    `json:"payload,omitempty"`
}

// Bytes converts OutboundMsg to Protobuf binary
func (m OutboundMsg) Bytes() []byte {
	sm := &ServerMsg{
		Type:   m.Type,
		RoomId: m.RoomID,
		GameId: m.GameID,
		Rev:    m.Rev,
		From:   m.From,
	}
	if m.Event != nil {
		pJSON, _ := json.Marshal(m.Event.Payload)
		sm.Event = &GameEvent{
			EventId:     m.Event.EventID,
			StreamId:    m.Event.StreamID,
			Seq:         m.Event.Seq,
			Type:        m.Event.Type,
			PayloadJson: string(pJSON),
		}
	}
	if m.State != nil {
		sJSON, _ := json.Marshal(m.State)
		sm.StateJson = string(sJSON)
	}
	if m.Payload != nil {
		pJSON, _ := json.Marshal(m.Payload)
		sm.PayloadJson = string(pJSON)
	}
	if m.Code != "" || m.Error != "" {
		sm.Error = &ErrorMsg{Code: m.Code, Error: m.Error}
	}

	b, err := proto.Marshal(sm)
	if err != nil {
		return []byte{}
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
	Payload json.RawMessage
}

// NewInboundMsg parses the binary frame into InboundMsg
func NewInboundMsg(data []byte) (InboundMsg, error) {
	var cm ClientMsg
	if err := proto.Unmarshal(data, &cm); err != nil {
		return InboundMsg{}, err
	}
	
	var payload json.RawMessage
	if cm.PayloadJson != "" {
		payload = json.RawMessage(cm.PayloadJson)
	}

	return InboundMsg{
		ID:      cm.Id,
		Type:    cm.Type,
		RoomID:  cm.RoomId,
		Rev:     uint64(cm.Rev),
		Payload: payload,
	}, nil
}

func (m InboundMsg) Command() (types.Command, error) {
	cmd := types.Command{
		ID:      m.ID,
		RoomID:  types.RoomID(m.RoomID),
		Type:    m.Type,
		Version: m.Rev,
	}
	if len(m.Payload) == 0 || string(m.Payload) == "null" {
		return cmd, nil
	}
	var p any
	if err := json.Unmarshal(m.Payload, &p); err != nil {
		return cmd, err
	}
	cmd.Payload = p
	return cmd, nil
}
