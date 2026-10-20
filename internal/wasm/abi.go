package wasm

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
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

func CallApply(ctx context.Context, mod api.Module, state, intent string, gameCtx GameContext) (string, error) {
	applyFn := mod.ExportedFunction("apply")
	if applyFn == nil {
		return "", fmt.Errorf("wasm: missing export %q", "apply")
	}
	allocFn := mod.ExportedFunction("allocate")
	if allocFn == nil {
		return "", fmt.Errorf("wasm: missing export %q", "allocate")
	}

	ctxB, err := json.Marshal(gameCtx)
	if err != nil { return "", err }

	writeMem := func(data []byte) (uint64, error) {
		res, err := allocFn.Call(ctx, api.EncodeI32(int32(len(data))))
		if err != nil { return 0, fmt.Errorf("allocate: %w", err) }
		if len(res) != 1 { return 0, fmt.Errorf("allocate: got %d results, want 1", len(res)) }
		ptr := uint32(api.DecodeI32(res[0]))
		if ptr == 0 { return 0, fmt.Errorf("allocate: guest returned 0") }
		if !mod.Memory().Write(ptr, data) { return 0, fmt.Errorf("memory write at %d failed", ptr) }
		return api.EncodeU32(ptr), nil
	}

	statePtr, err := writeMem([]byte(state)); if err != nil { return "", err }
	intentPtr, err := writeMem([]byte(intent)); if err != nil { return "", err }
	ctxPtr, err := writeMem(ctxB); if err != nil { return "", err }

	res, err := applyFn.Call(ctx,
		statePtr, api.EncodeI32(int32(len(state))),
		intentPtr, api.EncodeI32(int32(len(intent))),
		ctxPtr, api.EncodeI32(int32(len(ctxB))),
	)
	if err != nil { return "", err }
	if len(res) != 1 { return "", fmt.Errorf("apply: got %d results, want 1", len(res)) }
	outPtr := uint32(api.DecodeI32(res[0]))
	if outPtr == 0 { return "", fmt.Errorf("apply: guest returned 0") }

	hdr, ok := mod.Memory().Read(outPtr, 4)
	if !ok { return "", fmt.Errorf("read length failed") }
	n := binary.LittleEndian.Uint32(hdr)
	body, ok := mod.Memory().Read(outPtr+4, n)
	if !ok { return "", fmt.Errorf("read body failed") }
	return string(body), nil
}
