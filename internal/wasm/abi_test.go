package wasm

import (
	"context"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestCallApply(t *testing.T) {
	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)
	defer rt.Close(ctx)
	wasi_snapshot_preview1.Instantiate(ctx, rt)

	// Load the compiled WASM from games/test
	wasmBytes, err := os.ReadFile("../../games/test/game.wasm")
	if err != nil {
		t.Skip("game.wasm not built yet")
	}

	mod, err := rt.InstantiateWithConfig(ctx, wasmBytes, wazero.NewModuleConfig().WithStartFunctions())
	if err != nil {
		t.Fatal(err)
	}

	gctx := GameContext{
		GameID: "test",
		Rules:  GameRules{TurnTimeoutSec: 10},
		ActiveItems: []ShopItem{{ID: "item1", Type: "extra_time"}},
	}

	res, err := CallApply(ctx, mod, "init_state", "move", gctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Result: %s", res)
}
