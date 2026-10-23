package games

import "embed"

// 🛠️ جادوی Go: فایل‌های WASM را مستقیماً داخل خروجی کامپایل می‌کند
//go:embed */*.wasm
var WasmFS embed.FS
