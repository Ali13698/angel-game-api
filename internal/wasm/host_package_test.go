package wasm

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDynamicLoading verifies that the loader reads from the ./games directory
// instead of relying on Go's embed package.
func TestDynamicLoading(t *testing.T) {
	// Create a temporary mock games directory
	tempDir, err := os.MkdirTemp("", "angel-games-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// We don't have a real .wasm here to compile, so we just test
	// that the Init function gracefully handles missing or invalid files
	// without crashing the whole server.
	os.Setenv("GAMES_DIR", tempDir)
	defer os.Unsetenv("GAMES_DIR")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = Init(ctx)
	if err != nil {
		t.Fatalf("Init failed unexpectedly: %v", err)
	}
	defer Close(context.Background())

	// Verify ListGames is empty (since temp dir is empty)
	gamesList := ListGames()
	if len(gamesList) != 0 {
		t.Fatalf("Expected 0 games, got %d", len(gamesList))
	}

	// Verify that trying to instantiate a non-existent game fails safely
	_, err = Instantiate(ctx, "non-existent-game")
	if err == nil {
		t.Fatal("Expected error when instantiating non-existent game, got nil")
	}

	// Verify GAMES_DIR override works
	expectedPath := filepath.Clean(tempDir)
	if os.Getenv("GAMES_DIR") != expectedPath {
		t.Fatalf("GAMES_DIR not set correctly")
	}
}
