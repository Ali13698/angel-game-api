package games

import "embed"

//go:embed */*.wasm
var WasmFS embed.FS
