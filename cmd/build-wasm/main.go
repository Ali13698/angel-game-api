package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("🔨 Building WASM games...")
	
	entries, err := os.ReadDir("games")
	if err != nil {
		fmt.Println("No games directory found, skipping.")
		return
	}

	successCount := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		
		gameDir := filepath.Join("games", e.Name())
		wasmPath := filepath.Join(gameDir, "game.wasm")
		if _, err := os.Stat(filepath.Join(gameDir, "main.go")); err != nil {
			continue // فقط پوشه‌هایی که main.go دارند (بازی‌های واقعی)
		}

		fmt.Printf("➜ Compiling %s...\n", e.Name())
		
		// 🛠️ اصلاح قطعی دیپ‌سیک: حذف -buildmode=c-shared
		cmd := exec.Command("go", "build", "-o", "game.wasm", ".")
		cmd.Dir = gameDir
		cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0")
		
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error building %s:\n%s\n", e.Name(), string(out))
			os.Exit(1) // Fail the Docker build if a game fails
		} else {
			fmt.Printf("✅ Built %s successfully!\n", wasmPath)
			successCount++
		}
	}
	fmt.Printf("🎉 Done! %d games compiled.\n", successCount)
}
