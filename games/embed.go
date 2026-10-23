package games

import "embed"

// 🛠️ دقیقاً مسیر فایل‌ها را می‌دهیم تا اگر لودو ساخته نشد، بیلد با خطا متوقف شود
//go:embed ludo/game.wasm
//go:embed test/game.wasm
var WasmFS embed.FS
