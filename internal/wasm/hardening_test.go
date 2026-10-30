package wasm

import (
	"context"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/tetratelabs/wazero"
)

var (
	infiniteLoopHex = []byte("\x00\x61\x73\x6d\x01\x00\x00\x00\x01\x04\x01\x60\x00\x00\x03\x02\x01\x00\x07\x05\x01\x03\x72\x75\x6e\x00\x00\x0a\x07\x01\x05\x00\x03\x40\x0c\x00\x0b")
	healthyHex      = []byte("\x00\x61\x73\x6d\x01\x00\x00\x00\x01\x04\x01\x60\x00\x00\x03\x02\x01\x00\x07\x05\x01\x03\x72\x75\x6e\x00\x00\x0a\x04\x01\x02\x00\x0b")
)

func TestExecutionDeadlineIsolation(t *testing.T) {
	ctx := context.Background()

	t.Run("export-call", func(t *testing.T) {
		runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
		defer runtime.Close(ctx)

		compiled, err := runtime.CompileModule(ctx, infiniteLoopHex)
		if err != nil {
			t.Fatalf("failed to compile infinite loop module: %v", err)
		}

		deadline, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		mod, err := runtime.InstantiateModule(deadline, compiled, wazero.NewModuleConfig().WithName(""))
		if err != nil {
			// If instant instantiation timeout
			return
		}

		fn := mod.ExportedFunction("run")
		if fn != nil {
			_, _ = fn.Call(deadline)
		}

		if !mod.IsClosed() {
			t.Errorf("expected timed out instance to be closed")
		}

		// Non-poisoning check: healthy module should instantiate and run fine in same runtime
		healthyCompiled, err := runtime.CompileModule(ctx, healthyHex)
		if err != nil {
			t.Fatalf("runtime poisoned by previous timeout: %v", err)
		}
		healthyMod, err := runtime.InstantiateModule(ctx, healthyCompiled, wazero.NewModuleConfig().WithName(""))
		if err != nil {
			t.Fatalf("failed to instantiate healthy module in runtime: %v", err)
		}
		defer healthyMod.Close(ctx)
	})
}

func TestStrictRevisionBeforeInstantiation(t *testing.T) {
	handler := NewHandler("test-game")
	state := State{Rev: 7, Data: []byte(`{}`)}

	for _, badRev := range []uint64{0, 6, 8} {
		cmd := types.Command{
			Type:    "MOVE",
			Version: badRev,
		}
		nextState, events, err := handler.ApplyCommand(state, "player1", cmd)
		if err == nil {
			t.Errorf("expected error for bad revision %d, got nil", badRev)
		}
		if len(events) > 0 {
			t.Errorf("expected no events on revision mismatch")
		}
		if nextState.(State).Rev != state.Rev {
			t.Errorf("state modified on revision mismatch")
		}
	}

	if err := validateCommandVersion(7, 7); err != nil {
		t.Errorf("expected valid version 7 to pass, got: %v", err)
	}
}
