package main

import (
	"encoding/json"
	"unsafe"
)

type GameContext struct {
	GameID      string `json:"game_id"`
	Players     []any  `json:"players"`
	ActiveItems []any  `json:"active_items"`
	Rules       struct {
		TurnTimeoutSec int `json:"turn_timeout_sec"`
	} `json:"rules"`
}

//export allocate
func allocate(size int32) int32 {
	buf := make([]byte, size)
	return int32(uintptr(unsafe.Pointer(&buf[0])))
}

//export apply
func apply(statePtr int32, stateLen int32, intentPtr int32, intentLen int32, ctxPtr int32, ctxLen int32) int32 {
	// خواندن GameContext از حافظه
	ctxBytes := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ctxPtr))), ctxLen)
	var ctx GameContext
	json.Unmarshal(ctxBytes, &ctx)

	// منطق تستی: اگر آیتمی از فروشگاه داره، تایمر رو زیاد کن
	newState := "updated_state"
	if len(ctx.ActiveItems) > 0 && ctx.Rules.TurnTimeoutSec == 10 {
		newState = "state_with_extra_time"
	}

	outBytes := []byte(newState)
	outPtr := allocate(int32(len(outBytes)))
	copy(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(outPtr))), len(outBytes)), outBytes)
	return outPtr
}

func main() {}
