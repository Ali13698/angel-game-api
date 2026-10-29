package wasm

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/angel/angel/games"
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
	err := fs.WalkDir(games.WasmFS, ".", func(
		name string, entry fs.DirEntry, walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(name, ".wasm") {
			return nil
		}

		id := path.Dir(name)
		if id == "." {
			return nil
		}
		if _, exists := modules[id]; exists {
			return fmt.Errorf("wasm: duplicate game directory %q", id)
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() > maxModuleBytes {
			log.Printf("[wasm] skipped oversized module: %s", name)
			return nil
		}

		data, err := fs.ReadFile(games.WasmFS, name)
		if err != nil {
			return err
		}

		compileCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		compiled, err := runtime.CompileModule(compileCtx, data)
		cancel()
		if err != nil {
			log.Printf("[wasm] skipped %s: %v", name, err)
			return nil
		}

		modules[id] = compiled
		log.Printf("[wasm] compiled: %s", id)
		return nil
	})
	if err != nil {
		_ = runtime.Close(context.Background())
		return err
	}

	rt = runtime
	loaded = modules
	log.Printf("[wasm] compiled games: %d", len(loaded))
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
		return nil, fs.ErrNotExist
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
