package wasm

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const maxModuleBytes = 16 << 20

var (
	rt     wazero.Runtime
	mu     sync.RWMutex
	loaded map[string]wazero.CompiledModule
)

func runtimeConfig() wazero.RuntimeConfig {
	return wazero.NewRuntimeConfig().
		WithMemoryLimitPages(1024).
		WithCloseOnContextDone(true)
}

// Init now loads games dynamically from a directory on disk (default: "./games")
// This decouples the game engine from the core binary.
func Init(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if rt != nil {
		return fmt.Errorf("wasm: already initialized")
	}

	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig())
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		_ = runtime.Close(context.Background())
		return err
	}

	modules := make(map[string]wazero.CompiledModule)

	gamesDir := os.Getenv("GAMES_DIR")
	if gamesDir == "" {
		gamesDir = "./games" // Fallback to local directory
	}

	entries, err := os.ReadDir(gamesDir)
	if err != nil {
		log.Printf("[wasm] Warning: games directory %s not found. No games loaded.", gamesDir)
		rt = runtime
		loaded = modules
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		id := entry.Name()
		wasmPath := filepath.Join(gamesDir, id, "game.wasm")

		data, err := os.ReadFile(wasmPath)
		if err != nil {
			log.Printf("[wasm] skipped %s: missing game.wasm", id)
			continue
		}

		if len(data) > maxModuleBytes {
			log.Printf("[wasm] skipped oversized module: %s", id)
			continue
		}

		compileCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		compiled, err := runtime.CompileModule(compileCtx, data)
		cancel()
		if err != nil {
			log.Printf("[wasm] failed to compile %s: %v", id, err)
			continue
		}

		modules[id] = compiled
		log.Printf("[wasm] dynamically loaded: %s", id)
	}

	rt = runtime
	loaded = modules
	log.Printf("[wasm] total games loaded: %d", len(loaded))
	return nil
}

func ListGames() []string {
	mu.RLock()
	defer mu.RUnlock()

	ids := make([]string, 0, len(loaded))
	for id := range loaded {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func Instantiate(ctx context.Context, gameID string) (api.Module, error) {
	mu.RLock()
	runtime := rt
	compiled, exists := loaded[gameID]
	mu.RUnlock()

	if runtime == nil || !exists {
		return nil, fmt.Errorf("wasm: game not found: %s", gameID)
	}

	config := wazero.NewModuleConfig().
		WithName("").
		WithStartFunctions("_initialize")

	return runtime.InstantiateModule(ctx, compiled, config)
}

func Close(ctx context.Context) {
	mu.Lock()
	runtime := rt
	rt = nil
	loaded = nil
	mu.Unlock()

	if runtime != nil {
		_ = runtime.Close(ctx)
	}
}
