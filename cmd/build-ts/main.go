package main

import (
	"io"
	"log"
	"os"
	"github.com/evanw/esbuild/pkg/api"
)

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil { log.Fatal(err) }
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil { log.Fatal(err) }
	defer out.Close()
	io.Copy(out, in)
}

func main() {
	os.MkdirAll("internal/web/static", 0755)

	// 1. Build Profile
	res1 := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/profile/src/main.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(res1.Errors) > 0 { log.Fatalf("[build-ts] profile errors: %v", res1.Errors) }
	copyFile("web/profile/index.html", "internal/web/static/index.html")
	copyFile("web/profile/style.css", "internal/web/static/style.css")

	// 2. Build Lobby
	res2 := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/lobby/src/main.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/lobby.js",
	})
	if len(res2.Errors) > 0 { log.Fatalf("[build-ts] lobby errors: %v", res2.Errors) }
	copyFile("web/lobby/index.html", "internal/web/static/lobby.html")
	copyFile("web/lobby/style.css", "internal/web/static/lobby.css")

	log.Println("[build-ts] OK.")
}
