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

	cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", wasmPath, ".")
	cmd.Dir = filepath.Join("..", "..", "games", "test")
	cmd.Env = append(os.Environ(),
		"GOOS=wasip1",
		"GOARCH=wasm",
		"CGO_ENABLED=0", // این خط و os.Environ() مشکل رو حل می‌کنه
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s", err, out)
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

	// لاگ کردن لیست اکسپورت‌ها برای دیباگ
	t.Logf("wasm exports: %v", mod.ExportedFunctionNames())

	outStr, err := CallApply(ctx, mod, `{"state":"init"}`, `{"intent":"add_time"}`, GameContext{
		GameID: "test",
		Rules:  GameRules{TurnTimeoutSec: 30},
	})
	if err != nil {
		t.Fatalf("CallApply: %v", err)
	}
	t.Logf("output: %s", outStr)
	if outStr == "" {
		t.Fatal("empty output")
	}
}
