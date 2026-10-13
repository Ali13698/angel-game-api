package hub

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
)

// پین‌نوشت ۰.۱ — تأیید نگاشت کدهای خطا در لبهٔ wire.
//
// این تست ادعا می‌کند:
//   - ErrStaleVersion → CodeStaleVersion
//   - ErrWrongPhase   → CodeWrongPhase
//   - ErrBusy         → CodeBusy
//   - ErrClosed       → CodeServerClosed
//   - خطای نامشخص     → CodeRejected

func decodeFrame(t *testing.T, frame []byte) wire.OutboundMsg {
	t.Helper()
	var msg wire.OutboundMsg
	if err := json.Unmarshal(frame, &msg); err != nil {
		t.Fatalf("unmarshal frame: %v", err)
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
	frame := errFrameFor(room.ErrStaleVersion)
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeStaleVersion {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeStaleVersion)
	}
	if msg.Error == "" {
		t.Fatal("Error message is empty")
	}
}

func TestErrFrameFor_WrongPhase(t *testing.T) {
	frame := errFrameFor(room.ErrWrongPhase)
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeWrongPhase {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeWrongPhase)
	}
}

func TestErrFrameFor_Busy(t *testing.T) {
	frame := errFrameFor(room.ErrBusy)
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeBusy {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeBusy)
	}
}

func TestErrFrameFor_Closed(t *testing.T) {
	frame := errFrameFor(room.ErrClosed)
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeServerClosed {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeServerClosed)
	}
}

func TestErrFrameFor_Unknown(t *testing.T) {
	frame := errFrameFor(errors.New("some unclassified error"))
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeRejected {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeRejected)
	}
}

// تأیید می‌کند که errFrameFor با errors.Is روی خطاهای پیچیده هم کار می‌کند.
func TestErrFrameFor_Wrapped(t *testing.T) {
	wrapped := errors.Join(errors.New("context"), room.ErrStaleVersion)
	frame := errFrameFor(wrapped)
	msg := decodeFrame(t, frame)
	if msg.Code != wire.CodeStaleVersion {
		t.Fatalf("Code = %q, want %q", msg.Code, wire.CodeStaleVersion)
	}
}
