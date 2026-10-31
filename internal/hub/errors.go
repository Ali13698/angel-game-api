package hub

import (
	"github.com/angel/angel/internal/core/resource"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
)

func errFrameFor(err error) wire.OutboundMsg {
	switch err {
	case room.ErrStaleVersion:
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeStaleVersion, Error: "stale move"}
	case room.ErrWrongPhase:
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeWrongPhase, Error: "move rejected in current phase"}
	case room.ErrBusy:
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeBusy, Error: "room busy"}
	case room.ErrClosed:
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeServerClosed, Error: "room closed"}
	case resource.ErrCapacityFull:
		// 🔴 تبدیل خطای ظرفیت به پیام قابل فهم برای فرانت‌اند
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeCapacityFull, Error: "capacity full"}
	default:
		return wire.OutboundMsg{Type: wire.TypeError, Code: wire.CodeRejected, Error: "rejected"}
	}
}
