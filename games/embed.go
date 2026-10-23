package games

import "embed"

// 🛠️ به جای */*.wasm، دقیقاً مسیر فایل‌ها را می‌دهیم تا اگر لودو ساخته نشد، بیلد fail شود
//go:embed ludo/game.wasm
//go:embed test/game.wasm
var WasmFS embed.FS
