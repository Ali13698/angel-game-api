package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/tetratelabs/wazero/api"
)

type GameContext struct {
	GameID      string     `json:"game_id"`
	Players     []Player   `json:"players"`
	ActiveItems []ShopItem `json:"active_items"`
	Rules       GameRules  `json:"rules"`
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ShopItem struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Effect map[string]any `json:"effect"`
}

type GameRules struct {
	TurnTimeoutSec int `json:"turn_timeout_sec"`
}

func CallApply(ctx context.Context, mod api.Module, state string, intent string, gameCtx GameContext) (string, error) {
	applyFn := mod.ExportedFunction("apply")
	allocFn := mod.ExportedFunction("allocate")
	if applyFn == nil || allocFn == nil {
		return "", fmt.Errorf("wasm: missing apply or allocate function")
	}

	// Marshal inputs
	stateB := []byte(state)
	intentB := []byte(intent)
	ctxB, _ := json.Marshal(gameCtx)

	// Helper to write to guest memory
	writeMem := func(data []byte) uint64 {
		res, _ := allocFn.Call(ctx, uint64(len(data)))
		ptr := uint32(res[0])
		mod.Memory().Write(ptr, data)
		return uint64(ptr)
	}

	statePtr := writeMem(stateB)
	intentPtr := writeMem(intentB)
	ctxPtr := writeMem(ctxB)

	// Call apply(statePtr, stateLen, intentPtr, intentLen, ctxPtr, ctxLen)
	res, err := applyFn.Call(ctx, statePtr, uint64(len(stateB)), intentPtr, uint64(len(intentB)), ctxPtr, uint64(len(ctxB)))
	if err != nil {
		return "", err
	}

	// The function returns the pointer to the result string. (For simplicity in test, assuming fixed length or null-terminated)
	outPtr := uint32(res[0])
	outBytes, _ := mod.Memory().Read(outPtr, 64) // Read max 64 bytes for test
	return string(outBytes), nil
}
