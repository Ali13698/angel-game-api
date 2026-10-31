package hub

import (
	"errors"
	"github.com/angel/angel/internal/core/resource"
	"github.com/angel/angel/internal/room"
	"github.com/angel/angel/internal/wire"
)

var (
	errBadPlayRequest = errors.New("bad play request")
	errPlayNoMatch    = errors.New("no such match")
	errPlayNoRoom     = errors.New("room is full")
)

// errFrameFor returns []byte to match the signature in client_handle.go
func errFrameFor(err error) []byte {
	switch err {
	case room.ErrStaleVersion:
		return wire.ErrFrame(wire.CodeStaleVersion, "stale move")
	case room.ErrWrongPhase:
		return wire.ErrFrame(wire.CodeWrongPhase, "move rejected in current phase")
	case room.ErrBusy:
		return wire.ErrFrame(wire.CodeBusy, "room busy")
	case room.ErrClosed:
		return wire.ErrFrame(wire.CodeServerClosed, "room closed")
	case resource.ErrCapacityFull:
		return wire.ErrFrame(wire.CodeCapacityFull, "capacity full")
	default:
		return wire.ErrFrame(wire.CodeRejected, "rejected")
	}
}
