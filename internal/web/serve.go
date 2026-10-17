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
	
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// تنظیم MIME type برای فایل‌های JS
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		}
		http.FileServer(http.FS(sub)).ServeHTTP(w, r)
	})
	
	mux.Handle("/profile/", http.StripPrefix("/profile/", handler))
}
