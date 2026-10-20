package wasm

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestCallApply(t *testing.T) {
	ctx := context.Background()

	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("locate repo root: %v", err)
	}
	gameDir := filepath.Join(root, "games", "test")
	if _, err := os.Stat(filepath.Join(gameDir, "main.go")); err != nil {
		t.Fatalf("main.go missing at %s: %v", gameDir, err)
	}

	// Build into a temp dir; never read a checked-in game.wasm.
	wasmPath := filepath.Join(t.TempDir(), "game.wasm")
	cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", wasmPath, ".")
	cmd.Dir = gameDir
	cmd.Env = append(os.Environ(),
		"GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0",
		"GOFLAGS=", // ignore any GOFLAGS from the environment
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build games/test as wasip1 reactor: %v\n%s", err, out)
	}

	data, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("read built wasm: %v", err)
	}
	if len(data) < 8 || string(data[:4]) != "\x00asm" {
		t.Fatalf("built file is not a wasm module (%d bytes)", len(data))
	}
	t.Logf("built wasm: %d bytes", len(data))

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateWithConfig(ctx, data,
		wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil {
		t.Fatalf("instantiate: %v", err)
	}

	// Assert the exports actually exist and print them all if not.
	required := []string{"_initialize", "reset", "allocate", "apply"}
	for _, name := range required {
		if mod.ExportedFunction(name) == nil {
			t.Fatalf("export %q missing; exports present: %v", name, exportedNames(mod))
		}
	}

	out, err := CallApply(ctx, mod, `{"state":"init"}`, `{"intent":"add_time"}`, GameContext{
		GameID: "test",
		Rules:  GameRules{TurnTimeoutSec: 30},
	})
	if err != nil {
		t.Fatalf("CallApply: %v", err)
	}
	if out == "" {
		t.Fatal("CallApply returned empty result")
	}
	t.Logf("CallApply -> %s", out)
}

func exportedNames(mod api.Module) []string {
	defs := mod.ExportedFunctionDefinitions()
	names := make([]string, 0, len(defs))
	for n := range defs {
		names = append(names, n)
	}
	return names
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
