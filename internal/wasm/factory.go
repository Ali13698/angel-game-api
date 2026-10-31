package wasm

import (
	"context"
	"github.com/angel/angel/internal/core/types"
)

type Factory struct {
	GameID string
}

func NewFactory(gameID string) *Factory {
	return &Factory{GameID: gameID}
}

func (f *Factory) ID() types.GameID {
	return types.GameID(f.GameID)
}

func (f *Factory) New(ctx context.Context) (types.Game, error) {
	return NewHandler(ctx, f.GameID)
}
