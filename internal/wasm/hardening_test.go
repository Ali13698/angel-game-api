package wasm

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/tetratelabs/wazero"
)

const (
	infiniteLoopHex = "0061736d01000000010401600000030201000707010372756e00000a0901070003400c000b0b"
	healthyHex      = "0061736d01000000010401600000030201000707010372756e00000a040102000b"
)

func TestExecutionDeadlineIsolation(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "export-call"},
		{name: "start-function"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
			defer runtime.Close(ctx)

			infBytes, err := hex.DecodeString(infiniteLoopHex)
			if err != nil {
				t.Fatalf("decode loop hex: %v", err)
			}

			compiled, err := runtime.CompileModule(ctx, infBytes)
			if err != nil {
				t.Fatalf("compile loop module: %v", err)
			}

			deadlineCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer cancel()

			start := time.Now()
			mod, err := runtime.InstantiateModule(deadlineCtx, compiled, wazero.NewModuleConfig())
			if err == nil {
				fn := mod.ExportedFunction("run")
				if fn != nil {
					_, err = fn.Call(deadlineCtx)
				}
				_ = mod.Close(ctx)
			}
			elapsed := time.Since(start)

			if elapsed > 2*time.Second {
				t.Fatalf("execution took too long: %v", elapsed)
			}

			healthyBytes, err := hex.DecodeString(healthyHex)
			if err != nil {
				t.Fatalf("decode healthy hex: %v", err)
			}
			healthyCompiled, err := runtime.CompileModule(ctx, healthyBytes)
			if err != nil {
				t.Fatalf("compile healthy module after cancellation: %v", err)
			}
			healthyMod, err := runtime.InstantiateModule(ctx, healthyCompiled, wazero.NewModuleConfig())
			if err != nil {
				t.Fatalf("runtime corrupted after cancelled module: %v", err)
			}
			_ = healthyMod.Close(ctx)
		})
	}
}

func TestStrictRevisionBeforeInstantiation(t *testing.T) {
	currentRev := uint64(7)
	handler := NewHandler("not-loaded")
	st := &State{Rev: currentRev, Data: []byte(`{"turn":"A"}`)}

	for _, badRev := range []uint64{0, 6, 8} {
		cmd := types.Command{Type: "move", Version: badRev}
		if err := validateCommandVersion(cmd.Version, currentRev); err == nil {
			t.Fatalf("expected mismatch error for rev %d vs %d", badRev, currentRev)
		}
		_, _, err := handler.ApplyCommand(st, "user1", cmd)
		if err == nil {
			t.Fatalf("expected ApplyCommand to reject rev %d before instantiation", badRev)
		}
	}

	if err := validateCommandVersion(currentRev, currentRev); err != nil {
		t.Fatalf("expected valid revision match to pass: %v", err)
	}
}
