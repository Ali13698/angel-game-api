package main

import (
	"io"
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func copyIfExists(src, dst string) bool {
	in, err := os.Open(src)
	if err != nil { return false }
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil { return false }
	defer out.Close()
	io.Copy(out, in)
	return true
}

func buildFile(entry, outfile string) {
	res := api.Build(api.BuildOptions{
		EntryPoints:       []string{entry},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           outfile,
	})
	if len(res.Errors) > 0 {
		log.Fatalf("[build-ts] FATAL error building %s: %v", entry, res.Errors)
	}
	log.Printf("[build-ts] ✅ Built %s -> %s", entry, outfile)
}

func main() {
	if err := os.MkdirAll("internal/web/static", 0755); err != nil {
		log.Fatalf("[build-ts] mkdir: %v", err)
	}

	// بیلد پروفایل
	if _, err := os.Stat("web/profile/src/main.ts"); err == nil {
		buildFile("web/profile/src/main.ts", "internal/web/static/app.js")
		copyIfExists("web/profile/index.html", "internal/web/static/index.html")
		copyIfExists("web/profile/style.css", "internal/web/static/style.css")
	}

	// بیلد لابی (پنل بازی)
	if _, err := os.Stat("web/lobby/src/main.ts"); err == nil {
		buildFile("web/lobby/src/main.ts", "internal/web/static/lobby.js")
		copyIfExists("web/lobby/index.html", "internal/web/static/lobby.html")
		copyIfExists("web/lobby/style.css", "internal/web/static/lobby.css")
	} else {
		log.Println("[build-ts] ⚠️ web/lobby/src/main.ts not found. Skipping lobby build.")
	}

	log.Println("[build-ts] All frontend files processed.")
}
