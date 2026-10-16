package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:static
var staticFS embed.FS

func Register(mux *http.ServeMux) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil { return }
	fileServer := http.FileServer(http.FS(sub))
	mux.Handle("/profile/", http.StripPrefix("/profile/", fileServer))
}
