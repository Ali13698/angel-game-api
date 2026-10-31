package wasm

import (
	"context"
	"testing"
	"time"

	"github.com/angel/angel/internal/core/types" // 🔴 اضافه شد
)

func TestHardeningExecutionLimits(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	h, err := NewHandler(ctx, "test-game")
	if err != nil {
		t.Skip("test-game.wasm not found, skipping hardening test")
	}

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}

	state := &State{Rev: 10, Data: []byte("{}")}
	cmd := types.Command{
		Type:    "move",
		Version: 9,
	}

	_, _, err = h.ApplyCommand(state, "player1", cmd)
	if err == nil {
		t.Fatal("Expected stale version error, got nil")
	}
}
