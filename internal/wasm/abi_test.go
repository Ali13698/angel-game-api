package wasm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestCallApply(t *testing.T) {
	wasmPath := filepath.Join("..", "..", "games", "test", "game.wasm")
	data, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("read wasm: %v", err)
	}

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateWithConfig(ctx, data,
		wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil {
		t.Fatalf("instantiate: %v", err)
	}

	out, err := CallApply(ctx, mod, `{"state":"init"}`, `{"intent":"add_time"}`, GameContext{
		GameID: "test",
		Rules:  GameRules{TurnTimeoutSec: 30},
	})
	if err != nil {
		t.Fatalf("CallApply: %v", err)
	}
	t.Logf("output: %s", out)
	if out == "" {
		t.Fatal("empty output")
	}
}
