package wasm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/angel/angel/games"
)

func TestImportedHostPackage(t *testing.T) {
	manifestPath := os.Getenv("ANGEL_HOST_MANIFEST")
	if manifestPath == "" {
		t.Skip("run the Host Integration Check workflow to test imported games")
	}
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		SchemaVersion int `json:"schemaVersion"`
		Games []struct {
			ID string `json:"id"`
			WasmSha256 string `json:"wasmSha256"`
		} `json:"games"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != 1 || len(manifest.Games) == 0 {
		t.Fatal("invalid or empty imported manifest")
	}

	var expected []string
	for _, game := range manifest.Games {
		bytes, err := fs.ReadFile(games.WasmFS, game.ID+"/game.wasm")
		if err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(bytes)
		if hex.EncodeToString(hash[:]) != game.WasmSha256 {
			t.Fatalf("embedded WASM hash mismatch: %s", game.ID)
		}
		expected = append(expected, game.ID)
	}
	sort.Strings(expected)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer Close(context.Background())
	if got := ListGames(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("loaded games: %v; expected: %v", got, expected)
	}

	for _, id := range expected {
		t.Run(id, func(t *testing.T) {
			callCtx, stop := context.WithTimeout(ctx, 5*time.Second)
			defer stop()
			mod, err := Instantiate(callCtx, id)
			if err != nil {
				t.Fatal(err)
			}
			defer mod.Close(context.Background())
			if mod.Memory() == nil {
				t.Fatal("missing exported memory")
			}
			for _, name := range []string{"_initialize", "allocate", "reset", "apply"} {
				if mod.ExportedFunction(name) == nil {
					t.Fatalf("missing ABI export: %s", name)
				}
			}
		})
	}
}
