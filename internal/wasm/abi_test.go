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

// تابع کمکی برای چاپ لیست اکسپورت‌ها
func keys(m map[string]api.FunctionDefinition) []string {
	res := make([]string, 0, len(m))
	for k := range m { res = append(res, k) }
	return res
}

func TestCallApply(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil { t.Fatal(err) }
	
	gameDir := filepath.Join(root, "games", "test")
	if _, err := os.Stat(filepath.Join(gameDir, "main.go")); err != nil {
		t.Fatalf("main.go missing at %s: %v", gameDir, err)
	}

	wasmPath := filepath.Join(t.TempDir(), "game.wasm")
	cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", wasmPath, ".")
	cmd.Dir = gameDir
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0")
	
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	// لاگ کردن حجم و هدر فایل برای اطمینان از سالم بودن باینری
	fi, err := os.Stat(wasmPath)
	if err != nil { t.Fatalf("stat built wasm: %v", err) }
	head := make([]byte, 8)
	f, _ := os.Open(wasmPath)
	f.Read(head)
	f.Close()
	t.Logf("built wasm: size=%d magic=%q hex=% x", fi.Size(), head[:4], head[:8])

	data, err := os.ReadFile(wasmPath)
	if err != nil { t.Fatal(err) }

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.InstantiateWithConfig(ctx, data, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
	if err != nil { t.Fatal(err) }

	// چک کردن دقیق اکسپورت‌ها
	for _, n := range []string{"_initialize", "reset", "allocate", "apply"} {
		if fn := mod.ExportedFunction(n); fn == nil {
			t.Fatalf("export %q missing; have: %v", n, keys(mod.ExportedFunctionDefinitions()))
		}
	}

	outStr, err := CallApply(ctx, mod, `{"state":"init"}`, `{"intent":"add_time"}`, GameContext{
		GameID: "test",
		Rules:  GameRules{TurnTimeoutSec: 30},
	})
	if err != nil {
		t.Fatalf("CallApply: %v", err)
	}
	t.Logf("output: %s", outStr)
	if outStr == "" { t.Fatal("empty output") }
}
