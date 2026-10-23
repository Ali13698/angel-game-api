package wasm

import (
	"context"
	"io/fs"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/angel/angel/games"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
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
		return err
	}
	loaded = make(map[string]wazero.CompiledModule)

	// 🛠️ خواندن فایل‌ها مستقیماً از حافظه (Embedded FS)
	err := fs.WalkDir(games.WasmFS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(filePath, ".wasm") {
			return nil
		}

		// استخراج نام بازی از مسیر (مثلاً ludo/game.wasm -> ludo)
		gameID := path.Dir(filePath)
		if gameID == "." { return nil }

		data, err := fs.ReadFile(games.WasmFS, filePath)
		if err != nil {
			log.Printf("[wasm] Failed to read %s: %v", filePath, err)
			return nil
		}

		compiled, err := rt.CompileModule(ctx, data)
		if err != nil {
			log.Printf("[wasm] Failed to compile game %s: %v", gameID, err)
			return nil
		}

		mu.Lock()
		loaded[gameID] = compiled
		mu.Unlock()
		log.Printf("[wasm] ✅ Loaded embedded game: %s", gameID)
		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("[wasm] Total games embedded: %d", len(loaded))
	return nil
}

func ListGames() []string {
	mu.RLock()
	defer mu.RUnlock()
	keys := make([]string, 0, len(loaded))
	for k := range loaded {
		keys = append(keys, k)
	}
	return keys
}

func Instantiate(ctx context.Context, gameID string) (api.Module, error) {
	mu.RLock()
	compiled, ok := loaded[gameID]
	mu.RUnlock()
	if !ok {
		return nil, fs.ErrNotExist
	}
	return rt.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
}

func Close(ctx context.Context) {
	if rt != nil {
		_ = rt.Close(ctx)
	}
}
