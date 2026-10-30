package wasm

import (
	"encoding/json"
	"testing"

	"github.com/angel/angel/internal/core/types"
)

func TestStrictRevisionBeforeInstantiation(t *testing.T) {
	h := NewHandler("tic-tac-toe")
	s := State{Rev: 5, Data: json.RawMessage(`{}`)}

	// Version 0 or mismatched version must fail fast before calling WASM engine
	cmdWrong := types.Command{Type: "move", Version: 4}
	if _, _, err := h.ApplyCommand(s, "user1", cmdWrong); err == nil {
		t.Fatal("expected revision mismatch error, got nil")
	}

	cmdZero := types.Command{Type: "move", Version: 0}
	if _, _, err := h.ApplyCommand(s, "user1", cmdZero); err == nil {
		t.Fatal("expected revision mismatch error for version 0, got nil")
	}
}
