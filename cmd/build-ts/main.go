package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	log.Println("[tsc] compiling TypeScript to JavaScript using Go esbuild...")
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{"web/profile/app.ts"},
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
		Outfile:           "internal/web/static/app.js",
	})
	if len(result.Errors) > 0 {
		log.Printf("[tsc] errors: %v", result.Errors)
		os.Exit(1)
	}
	log.Println("[tsc] OK: web/profile/app.ts -> internal/web/static/app.js")
}
