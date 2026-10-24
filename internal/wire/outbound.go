package wire

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
)

type OutboundMsg struct {
	Type    string      `json:"type"`
	Code    string      `json:"code,omitempty"`
	RoomID  string      `json:"room_id,omitempty"`
	GameID  string      `json:"game_id,omitempty"`
	Rev     uint64      `json:"rev,omitempty"`
	Events  []*GameEvent `json:"events,omitempty"` // ✅ تغییر به آرایه
	State   any         `json:"state,omitempty"`
	From    string      `json:"from,omitempty"`
	Error   string      `json:"error,omitempty"`
	Payload any         `json:"payload,omitempty"`
}

func (m OutboundMsg) Bytes() []byte {
	sm := &ServerMsg{Type: m.Type, RoomId: m.RoomID, GameId: m.GameID, Rev: int64(m.Rev), From: m.From}
	
	// ✅ تغییر در اختصاص دادن رویدادها
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
	if err != nil {
		return []byte{}
	}
	return b
}
