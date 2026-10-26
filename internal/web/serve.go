package web

import (
	"embed"
	"net/http"
	"path"
	"strconv"
	"strings"
)

//go:embed all:static
var staticFS embed.FS

const noCacheHeader = "no-store, no-cache, must-revalidate, max-age=0"

func contentTypeOf(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js", ".mjs":
		return "application/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".proto":
		return "text/plain; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	case ".wasm":
		return "application/wasm"
	}
	return "application/octet-stream"
}

func serveStatic(w http.ResponseWriter, r *http.Request, name string) bool {
	name = strings.TrimPrefix(name, "/")
	if name == "" || name == "/" {
		name = "index.html"
	}
	
	// اگر مسیر با / تمام شده بود، فایل index.html داخل آن پوشه را سرو کن
	if strings.HasSuffix(name, "/") {
		name += "index.html"
	}
	
	if strings.Contains(name, "..") {
		return false
	}

	data, err := staticFS.ReadFile("static/" + name)
	if err != nil {
		return false
	}

	h := w.Header()
	h.Set("Cache-Control", noCacheHeader)
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
	h.Set("Content-Type", contentTypeOf(name))
	h.Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(data)
	}
	return true
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		
		// ۱. تلاش برای سرو فایل دقیق (مثلا /games/ludo/network.js)
		if serveStatic(w, r, name) {
			return
		}
		
		// ۲. اگر مسیر پسوند نداشت، احتمالا یک پوشه است (مثلا /games/ludo)
		if path.Ext(name) == "" {
			if serveStatic(w, r, name+"/index.html") {
				return
			}
		}
		
		// ۳. SPA Fallback (اگر هیچ فایلی پیدا نشد، صفحه اصلی پروفایل را بده)
		if serveStatic(w, r, "index.html") {
			return
		}
		
		http.NotFound(w, r)
	})
}
