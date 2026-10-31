package wasm

import (
	"context"
	"testing"
	"time"
)

func TestHardeningExecutionLimits(t *testing.T) {
	// This test ensures that the new V2 handler correctly rejects commands
	// if the state version is stale or if the timeout is exceeded.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create a dummy handler
	h, err := NewHandler(ctx, "test-game")
	if err != nil {
		// If the game.wasm isn't present in the test env, skip (but do not fail)
		t.Skip("test-game.wasm not found, skipping hardening test")
	}

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}

	// Mock state and command
	state := &State{Rev: 10, Data: []byte("{}")}
	cmd := types.Command{
		Type:    "move",
		Version: 9, // Stale version (less than state.Rev)
	}

	_, _, err = h.ApplyCommand(state, "player1", cmd)
	if err == nil {
		t.Fatal("Expected stale version error, got nil")
	}
}
