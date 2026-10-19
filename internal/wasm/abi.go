package wasm

import (
	"encoding/json"
	"fmt"
)

// GameContext تمام دیتای تزریقی از هسته به بازی را نگه می‌دارد
type GameContext struct {
	GameID      string           `json:"game_id"`
	Players     []Player         `json:"players"`
	ActiveItems []ShopItem       `json:"active_items"`
	Rules       GameRules        `json:"rules"`
}

type Player struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	IsBot  bool   `json:"is_bot"`
}

type ShopItem struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"` // e.g., "dice_multiplier", "extra_time"
	Effect map[string]any `json:"effect"`
}

type GameRules struct {
	TurnTimeoutSec int `json:"turn_timeout_sec"` // مثلا 10 ثانیه
	MaxPlayers     int `json:"max_players"`
}

// GameIntent قصد بازیکن برای انجام حرکت
type GameIntent struct {
	Action  string         `json:"action"` // e.g., "roll", "move"
	Payload map[string]any `json:"payload"`
}

// GameResult خروجی بازی پس از حرکت
type GameResult struct {
	NewState string   `json:"new_state"`
	Events   []string `json:"events"`
}

// BuildInitPayload ساخت JSON برای ارسال به بازی
func BuildInitPayload(ctx GameContext) string {
	b, _ := json.Marshal(ctx)
	return string(b)
}

// MockApply (شبیه‌سازی تست اتصال)
func MockApply(state string, intent string) (string, error) {
	// اینجا در آینده به حافظه WASM وصل می‌شویم
	return fmt.Sprintf("State: %s | Intent: %s", state, intent), nil
}
