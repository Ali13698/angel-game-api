package wasm

import (
	"context"

	"github.com/angel/angel/internal/core/types"
)

// Factory creates isolated WASM instances for each room.
type Factory struct {
	GameID string
}

func NewFactory(gameID string) *Factory {
	return &Factory{GameID: gameID}
}

func (f *Factory) ID() types.GameID {
	return types.GameID(f.GameID)
}

// New creates a fresh WASM handler (a new isolated module) for a specific room.
func (f *Factory) New(ctx context.Context) (types.Game, error) {
	return NewHandler(ctx, f.GameID)
}

var _ types.Game = (*WasmHandler)(nil)
