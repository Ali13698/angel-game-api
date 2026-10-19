package wasm

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var (
	rt     wazero.Runtime
	mu     sync.RWMutex
	loaded map[string]wazero.CompiledModule
)

func Init(ctx context.Context) error {
	rt = wazero.NewRuntime(ctx)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		return fmt.Errorf("wasm: failed to instantiate WASI: %w", err)
	}
	loaded = make(map[string]wazero.CompiledModule)
	gamesDir := "./games"
	os.MkdirAll(gamesDir, 0755)
	entries, _ := os.ReadDir(gamesDir)

	for _, entry := range entries {
		if !entry.IsDir() { continue }
		gameID := entry.Name()
		wasmPath := filepath.Join(gamesDir, gameID, "game.wasm")
		if _, err := os.Stat(wasmPath); os.IsNotExist(err) { continue }

		data, _ := os.ReadFile(wasmPath)
		compiled, err := rt.CompileModule(ctx, data)
		if err != nil {
			log.Printf("[wasm] Failed to compile game %s: %v", gameID, err)
			continue
		}
		mu.Lock()
		loaded[gameID] = compiled
		mu.Unlock()
		log.Printf("[wasm] Loaded game: %s", gameID)
	}
	log.Printf("[wasm] OK. Total games: %d", len(loaded))
	return nil
}

func GetGame(gameID string) (wazero.CompiledModule, bool) {
	mu.RLock()
	defer mu.RUnlock()
	mod, ok := loaded[gameID]
	return mod, ok
}

func Close(ctx context.Context) {
	if rt != nil { _ = rt.Close(ctx) }
}
