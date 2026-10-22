package items

import (
	"github.com/angel/angel/internal/wasm"
)

// MapToWasmItems وظیفه تبدیل آیتم‌های دیتابیس به فرمتی که بازی WASM می‌فهمد را دارد
func MapToWasmItems(inv []InfiniteItem) []wasm.ShopItem {
	if len(inv) == 0 {
		return nil
	}

	out := make([]wasm.ShopItem, 0, len(inv))
	for _, it := range inv {
		out = append(out, wasm.ShopItem{
			ID:     it.ID,
			Type:   string(it.Type),
			Effect: it.Meta,
		})
	}
	return out
}
