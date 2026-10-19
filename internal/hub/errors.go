package hub
import ("errors";"github.com/angel/angel/internal/room";"github.com/angel/angel/internal/wire")
func errFrameFor(err error) []byte {
	switch {
	case errors.Is(err, room.ErrStaleVersion):
		return wire.ErrFrame(wire.CodeStaleVersion, err.Error())
	case errors.Is(err, room.ErrWrongPhase):
		return wire.ErrFrame(wire.CodeWrongPhase, err.Error())
	case errors.Is(err, room.ErrBusy):
		return wire.ErrFrame(wire.CodeBusy, err.Error())
	case errors.Is(err, room.ErrClosed):
		return wire.ErrFrame(wire.CodeServerClosed, err.Error())
	default:
		return wire.ErrFrame(wire.CodeRejected, err.Error())
	}
}
