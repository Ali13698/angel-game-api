package wasm

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

const maxABIPayload = 256 << 10

type GameContext struct {
	GameID      string    `json:"game_id"`
	Players     []Player  `json:"players"`
	ActiveItems []any     `json:"active_items"`
	Rules       GameRules `json:"rules"`
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

func CallApply(
	ctx context.Context,
	mod api.Module,
	state, intent string,
	gameCtx GameContext,
) (string, error) {
	apply := mod.ExportedFunction("apply")
	allocate := mod.ExportedFunction("allocate")
	memory := mod.Memory()

	if apply == nil || allocate == nil || memory == nil {
		return "", fmt.Errorf("wasm: missing apply, allocate or memory")
	}

	contextBytes, err := json.Marshal(gameCtx)
	if err != nil {
		return "", err
	}

	inputs := [][]byte{
		[]byte(state),
		[]byte(intent),
		contextBytes,
	}
	total := 0
	for _, input := range inputs {
		total += len(input)
	}
	if total > maxABIPayload {
		return "", fmt.Errorf("wasm: input exceeds ABI limit")
	}

	write := func(data []byte) (uint64, error) {
		result, err := allocate.Call(ctx, uint64(len(data)))
		if err != nil {
			return 0, fmt.Errorf("wasm allocate: %w", err)
		}
		if len(result) != 1 || result[0] == 0 {
			return 0, fmt.Errorf("wasm: invalid allocation result")
		}
		ptr := uint32(result[0])
		if !memory.Write(ptr, data) {
			return 0, fmt.Errorf("wasm: allocation outside memory")
		}
		return uint64(ptr), nil
	}

	args := make([]uint64, 0, 6)
	for _, input := range inputs {
		ptr, err := write(input)
		if err != nil {
			return "", err
		}
		args = append(args, ptr, uint64(len(input)))
	}

	result, err := apply.Call(ctx, args...)
	if err != nil {
		return "", fmt.Errorf("wasm apply: %w", err)
	}
	if len(result) != 1 || result[0] == 0 {
		return "", fmt.Errorf("wasm: command rejected or invalid result")
	}

	ptr := uint32(result[0])
	if ptr > ^uint32(0)-4 {
		return "", fmt.Errorf("wasm: output pointer overflow")
	}

	header, ok := memory.Read(ptr, 4)
	if !ok {
		return "", fmt.Errorf("wasm: invalid output header")
	}

	size := binary.LittleEndian.Uint32(header)
	if size == 0 || size > maxABIPayload {
		return "", fmt.Errorf("wasm: output exceeds ABI limit")
	}

	body, ok := memory.Read(ptr+4, size)
	if !ok {
		return "", fmt.Errorf("wasm: output outside memory")
	}
	return string(body), nil
}
