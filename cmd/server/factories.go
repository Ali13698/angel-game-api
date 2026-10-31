package main

import (
	"context"
	"log"

	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/games/echo"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/wasm"
)

type echoFactory struct{}

func (echoFactory) ID() types.GameID { return types.GameID("echo") }

func (echoFactory) New(_ context.Context) (types.Game, error) {
	return echo.New(), nil
}

func registerGames(ctx context.Context, reg *registry.Registry) (int, error) {
	if err := reg.Register(echoFactory{}); err != nil {
		return 0, err
	}
	log.Printf("[startup] registered game: echo")

	if err := wasm.Init(ctx); err != nil {
		return 0, err
	}

	registered := 0
	for _, gameID := range wasm.ListGames() {
		if err := reg.Register(wasm.NewFactory(gameID)); err != nil {
			log.Printf("[wasm] registration failed %s: %v", gameID, err)
			continue
		}
		registered++
		log.Printf("[wasm] registered: %s", gameID)
	}
	log.Printf("[wasm] registered games: %d", registered)
	return registered, nil
}
