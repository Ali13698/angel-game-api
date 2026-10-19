package wasm

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var ( rt wazero.Runtime; mu sync.RWMutex; loaded map[string]wazero.CompiledModule )

func Init(ctx context.Context) error {
	rt = wazero.NewRuntime(ctx)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil { return fmt.Errorf("wasm: failed to instantiate WASI: %w", err) }
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
		if err != nil { log.Printf("[wasm] Failed to compile game %s: %v", gameID, err); continue }
		mu.Lock(); loaded[gameID] = compiled; mu.Unlock()
		log.Printf("[wasm] Loaded game: %s", gameID)
	}
	log.Printf("[wasm] OK. Total games: %d", len(loaded))
	return nil
}

func ListGames() []string {
	mu.RLock(); defer mu.RUnlock()
	keys := make([]string, 0, len(loaded))
	for k := range loaded { keys = append(keys, k) }
	return keys
}

func Instantiate(ctx context.Context, gameID string) (api.Module, error) {
	mu.RLock(); compiled, ok := loaded[gameID]; mu.RUnlock()
	if !ok { return nil, fmt.Errorf("wasm: game %s not compiled", gameID) }
	return rt.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
}

func Close(ctx context.Context) { if rt != nil { _ = rt.Close(ctx) } }
