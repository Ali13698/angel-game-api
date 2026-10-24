package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("🔨 Building WASM games as Reactors...")
	entries, err := os.ReadDir("games")
	if err != nil { return }

	for _, e := range entries {
		if !e.IsDir() { continue }
		gameDir := filepath.Join("games", e.Name())
		if _, err := os.Stat(filepath.Join(gameDir, "main.go")); err != nil { continue }

		fmt.Printf("➜ Compiling %s...\n", e.Name())
		// ✅ استفاده از buildmode=c-shared برای ساخت ری‌اکتور و رفع خطای notInitialized
		cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", "game.wasm", ".")
		cmd.Dir = gameDir
		cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0")
		
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error building %s:\n%s\n", e.Name(), string(out))
			os.Exit(1)
		}
		fmt.Printf("✅ Built %s successfully!\n", e.Name())
	}
}
