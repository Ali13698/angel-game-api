package main

import (
	"io"
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func copyFile(src, dst string) {
	in, err := os.Open(src); if err != nil { log.Fatal(err) }; defer in.Close()
	out, err := os.Create(dst); if err != nil { log.Fatal(err) }; defer out.Close()
	io.Copy(out, in)
}

func main() {
	os.MkdirAll("internal/web/static", 0755)
	os.Remove("internal/web/static/app.js")

	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/profile/src/main.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(result.Errors) > 0 { log.Fatalf("[build-ts] errors: %v", result.Errors) }
	
	copyFile("web/profile/index.html", "internal/web/static/index.html")
	copyFile("web/profile/style.css", "internal/web/static/style.css")
	log.Println("[build-ts] OK.")
}
