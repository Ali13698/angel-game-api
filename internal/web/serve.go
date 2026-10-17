package web

import (
	"embed"
	"net/http"
	"strings"
)

//go:embed all:static
var staticFS embed.FS

func Register(mux *http.ServeMux) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		
		path := strings.TrimPrefix(r.URL.Path, "/profile/")
		if path == "" || path == "/" {
			path = "index.html"
		}
		
		if idx := strings.Index(path, "?"); idx != -1 {
			path = path[:idx]
		}

		data, err := staticFS.ReadFile("static/" + path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}

		w.Write(data)
	})
	mux.Handle("/profile/", handler)
}
