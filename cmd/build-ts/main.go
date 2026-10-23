package main

import (
	"io"
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func copyIfExists(src, dst string) bool {
	in, err := os.Open(src)
	if err != nil {
		log.Printf("[build-ts] SKIP (missing source): %s", src)
		return false
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Printf("[build-ts] SKIP (cannot create %s): %v", dst, err)
		return false
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		log.Printf("[build-ts] copy error %s -> %s: %v", src, dst, err)
		return false
	}
	log.Printf("[build-ts] copied %s -> %s", src, dst)
	return true
}

func exists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func main() {
	if err := os.MkdirAll("internal/web/static", 0755); err != nil {
		log.Fatalf("[build-ts] mkdir: %v", err)
	}

	if !exists("web/profile/src/main.ts") {
		log.Fatalf("[build-ts] FATAL: web/profile/src/main.ts not found")
	}

	// ۱. بیلد پروفایل
	res1 := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/profile/src/main.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(res1.Errors) > 0 {
		log.Fatalf("[build-ts] profile errors: %v", res1.Errors)
	}
	log.Println("[build-ts] app.js built")

	// کپی استایل و HTML (از پوشه profile)
	copyIfExists("web/profile/index.html", "internal/web/static/index.html")
	copyIfExists("web/profile/style.css", "internal/web/static/style.css")

	// ۲. بیلد لابی (اگر وجود داشته باشد)
	if exists("web/lobby/src/main.ts") {
		res2 := api.Build(api.BuildOptions{
			EntryPoints:       []string{"web/lobby/src/main.ts"},
			Bundle:            true,
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
			Write:             true,
			Outfile:           "internal/web/static/lobby.js",
		})
		if len(res2.Errors) > 0 {
			log.Fatalf("[build-ts] lobby errors: %v", res2.Errors)
		}
		log.Println("[build-ts] lobby.js built")

		copyIfExists("web/lobby/style.css", "internal/web/static/lobby.css")
	} else {
		log.Println("[build-ts] lobby source not found — skipping lobby build")
	}

	log.Println("[build-ts] OK.")
}
