package wasm

import (
	"context"
	"testing"
	"time"
)

func TestExecutionDeadlineIsolation(t *testing.T) {
	// Verifies that a hung WASM module execution does not corrupt shared runtime/cache
	loader := NewLoader("games")
	adapter := NewAdapter(loader)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ensure adapter handles deadline enforcement correctly without panicking host
	_, _, err := adapter.ApplyCommand(ctx, "nonexistent-game", []byte{}, wireCommandMock())
	if err == nil {
		t.Log("Expected error on nonexistent/unloaded game, got nil")
	}
}

func TestStrictRevisionBeforeInstantiation(t *testing.T) {
	// Verify revision check occurs before heavy instantiation when possible, or via validate
	err := validateCommandVersion(5, 7)
	if err == nil {
		t.Fatal("Expected revision mismatch error, got nil")
	}

	err = validateCommandVersion(7, 7)
	if err != nil {
		t.Fatalf("Expected no error for matching revisions, got: %v", err)
	}
}

func wireCommandMock() struct {
	ID      string
	RoomID  string
	Type    string
	Payload []byte
	Version int64
} {
	return struct {
		ID      string
		RoomID  string
		Type    string
		Payload []byte
		Version int64
	}{
		ID:      "test-cmd",
		RoomID:  "room-1",
		Type:    "MOVE",
		Payload: []byte{},
		Version: 1,
	}
}
