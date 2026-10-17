package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:static
var staticFS embed.FS

func Register(mux *http.ServeMux) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil { return }
	fileServer := http.FileServer(http.FS(sub))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("/profile/", http.StripPrefix("/profile/", handler))
}
