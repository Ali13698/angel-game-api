package wasm

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestCallApply(t *testing.T) {
	tmpDir := t.TempDir()
	wasmPath := filepath.Join(tmpDir, "game.wasm")

	// حذف فلگ اشتباه c-shared که برای wasip1 پشتیبانی نمی‌شود
	cmd := exec.Command("go", "build", "-o", wasmPath, ".")
	cmd.Dir = filepath.Join("..", "..", "games", "test")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build wasm: %v\nOutput: %s", err, out)
	}

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
