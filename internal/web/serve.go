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
	
	// سرو کردن فایل‌ها با http.FileServer
	fileServer := http.FileServer(http.FS(sub))
	
	// برای اینکه وقتی /profile/ باز شد، index.html لود بشه
	mux.Handle("/profile/", http.StripPrefix("/profile/", fileServer))
}
