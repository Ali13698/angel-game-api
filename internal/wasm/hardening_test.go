package wasm

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types"
	"github.com/tetratelabs/wazero"
)

// (module (func (export "run") (loop br 0)))
const infiniteLoopHex = "0061736d01000000010401600000030201000707010372756e00000a0901070003400c000b0b"

// Same function, but it returns immediately.
const healthyHex = "0061736d01000000010401600000030201000707010372756e00000a040102000b"

func TestExecutionDeadlineIsolation(t *testing.T) {
	for _, startup := range []bool{false, true} {
		name := "export-call"
		if startup {
			name = "start-function"
		}
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
			defer runtime.Close(ctx)

			bytes, err := hex.DecodeString(infiniteLoopHex)
			if err != nil {
				t.Fatal(err)
			}
			compiled, err := runtime.CompileModule(ctx, bytes)
			if err != nil {
				t.Fatal(err)
			}
			defer compiled.Close(ctx)

			config := wazero.NewModuleConfig().WithName("").WithStartFunctions()
			if startup {
				config = config.WithStartFunctions("run")
			}

			deadline, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer cancel()

			started := time.Now()
			mod, err := runtime.InstantiateModule(deadline, compiled, config)
			if !startup {
				if err != nil {
					t.Fatal(err)
				}
				defer mod.Close(ctx)
				_, err = mod.ExportedFunction("run").Call(deadline)
				if !mod.IsClosed() {
					t.Fatal("timed-out instance must be closed")
				}
			}

			if err == nil || deadline.Err() != context.DeadlineExceeded {
				t.Fatalf("expected deadline interruption, got %v / %v", err, deadline.Err())
			}

			if time.Since(started) > 2*time.Second {
				t.Fatal("interruption exceeded tolerance")
			}

			// A failed instance must not poison the shared runtime.
			healthyBytes, err := hex.DecodeString(healthyHex)
			if err != nil {
				t.Fatal(err)
			}
			healthy, err := runtime.InstantiateWithConfig(ctx, healthyBytes,
				wazero.NewModuleConfig().WithName("").WithStartFunctions())
			if err != nil {
				t.Fatal(err)
			}
			defer healthy.Close(ctx)

			if _, err := healthy.ExportedFunction("run").Call(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStrictRevisionBeforeInstantiation(t *testing.T) {
	for _, revision := range []uint64{0, 6, 8} {
		if validateCommandVersion(revision, 7) == nil {
			t.Fatalf("accepted %d", revision)
		}
		state := &State{Rev: 7}
		next, events, err := NewHandler("not-loaded").ApplyCommand(state, "user", types.Command{Version: revision})
		if err == nil || next != state || len(events) != 0 {
			t.Fatal("invalid command changed state")
		}
	}
	if err := validateCommandVersion(7, 7); err != nil {
		t.Fatal(err)
	}
}
