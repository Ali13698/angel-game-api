package hub

import (
	"errors"
	"testing"

	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
	"google.golang.org/protobuf/proto"
)

func decodeFrame(t *testing.T, frame []byte) wire.OutboundMsg {
	t.Helper()
	var sm wire.ServerMsg
	if err := proto.Unmarshal(frame, &sm); err != nil {
		t.Fatalf("unmarshal frame: %v", err)
	}
	msg := wire.OutboundMsg{
		Type:   sm.GetType(),
		RoomID: sm.GetRoomId(),
		GameID: sm.GetGameId(),
		Rev:    uint64(sm.GetRev()),
		From:   sm.GetFrom(),
	}
	if e := sm.GetError(); e != nil {
		msg.Code = e.GetCode()
		msg.Error = e.GetError()
	}
	if msg.Type != wire.TypeError {
		t.Fatalf("Type = %q, want %q", msg.Type, wire.TypeError)
	}
	if msg.Code == "" {
		t.Fatal("Code is empty")
	}
	return msg
}

func TestErrFrameFor_StaleVersion(t *testing.T) {
	msg := decodeFrame(t, errFrameFor(room.ErrStaleVersion))
	if msg.Code != wire.CodeStaleVersion {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeStaleVersion)
	}
	if msg.Error == "" {
		t.Fatal("Error message is empty")
	}
}

func TestErrFrameFor_WrongPhase(t *testing.T) {
	msg := decodeFrame(t, errFrameFor(room.ErrWrongPhase))
	if msg.Code != wire.CodeWrongPhase {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeWrongPhase)
	}
}

func TestErrFrameFor_Busy(t *testing.T) {
	msg := decodeFrame(t, errFrameFor(room.ErrBusy))
	if msg.Code != wire.CodeBusy {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeBusy)
	}
}

func TestErrFrameFor_Closed(t *testing.T) {
	msg := decodeFrame(t, errFrameFor(room.ErrClosed))
	if msg.Code != wire.CodeServerClosed {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeServerClosed)
	}
}

func TestErrFrameFor_Unknown(t *testing.T) {
	msg := decodeFrame(t, errFrameFor(errors.New("some unclassified error")))
	if msg.Code != wire.CodeRejected {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeRejected)
	}
}

func TestErrFrameFor_Wrapped(t *testing.T) {
	wrapped := errors.Join(errors.New("context"), room.ErrStaleVersion)
	msg := decodeFrame(t, errFrameFor(wrapped))
	if msg.Code != wire.CodeStaleVersion {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeStaleVersion)
	}
}
