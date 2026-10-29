package wasm

import (
	"context"
	"testing"
	"time"

	"github.com/tetratelabs/wazero"
)

// A minimal module exporting spin(), which loops forever.
var infiniteLoopModule = []byte{
	0, 97, 115, 109, 1, 0, 0, 0,
	1, 4, 1, 96, 0, 0,
	3, 2, 1, 0,
	7, 8, 1, 4, 115, 112, 105, 110, 0, 0,
	10, 9, 1, 7, 0, 3, 64, 12, 0, 11, 11,
}

func TestExecutionDeadlineClosesGuest(t *testing.T) {
	ctx := context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
	defer runtime.Close(ctx)

	mod, err := runtime.InstantiateWithConfig(
		ctx,
		infiniteLoopModule,
		wazero.NewModuleConfig().WithName("").WithStartFunctions(),
	)
	if err != nil {
		t.Fatal(err)
	}

	callCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := mod.ExportedFunction("spin").Call(callCtx)
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("infinite loop returned without an error")
		}
		if !mod.IsClosed() {
			t.Fatal("timed-out guest was not closed")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("execution deadline did not stop the guest")
	}
}

func TestAnonymousInstancesDoNotConflict(t *testing.T) {
	ctx := context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
	defer runtime.Close(ctx)

	compiled, err := runtime.CompileModule(ctx, infiniteLoopModule)
	if err != nil {
		t.Fatal(err)
	}

	config := wazero.NewModuleConfig().
		WithName("").
		WithStartFunctions()

	first, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close(ctx)

	second, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		t.Fatalf("second instance conflicted: %v", err)
	}
	defer second.Close(ctx)
}

func TestFailedGuestDoesNotCloseAnotherGuest(t *testing.T) {
	ctx := context.Background()
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
	defer runtime.Close(ctx)

	compiled, err := runtime.CompileModule(ctx, infiniteLoopModule)
	if err != nil {
		t.Fatal(err)
	}

	config := wazero.NewModuleConfig().
		WithName("").
		WithStartFunctions()

	first, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		t.Fatal(err)
	}
	second, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close(ctx)

	callCtx, cancel := context.WithCancel(ctx)
	cancel()

	if _, err := first.ExportedFunction("spin").Call(callCtx); err == nil {
		t.Fatal("cancelled execution unexpectedly succeeded")
	}
	if second.IsClosed() {
		t.Fatal("cancelling one guest closed another")
	}
}
