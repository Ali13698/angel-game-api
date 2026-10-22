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
		
		fmt.Printf("➜ Compiling %s...\n", e.Name())
		
		cmd := exec.Command("go", "build", "-buildmode=c-shared", "-o", "game.wasm", ".")
		cmd.Dir = gameDir
		cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm", "CGO_ENABLED=0")
		
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Error building %s:\n%s\n", e.Name(), string(out))
		} else {
			fmt.Printf("✅ Built %s successfully!\n", wasmPath)
			successCount++
		}
	}
	
	fmt.Printf("🎉 Done! %d games compiled.\n", successCount)
}
