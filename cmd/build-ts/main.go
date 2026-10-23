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

func exists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func main() {
	if err := os.MkdirAll("internal/web/static", 0755); err != nil {
		log.Fatalf("[build-ts] mkdir: %v", err)
	}

	// بیلد پروفایل
	res := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/profile/src/main.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(res.Errors) > 0 {
		log.Fatalf("[build-ts] profile errors: %v", res.Errors)
	}
	log.Println("[build-ts] app.js built successfully")

	copyIfExists("web/profile/index.html", "internal/web/static/index.html")
	copyIfExists("web/profile/style.css", "internal/web/static/style.css")
}
