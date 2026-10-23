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
	// 🛠️ حذف شرط HasSuffix(name, "/") تا مسیرهای پایانی با اسلش هم بررسی شوند
	name = strings.TrimPrefix(name, "/")
	if name == "" {
		return false
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
		reqPath := r.URL.Path
		// 🛠️ اگر مسیر به اسلش ختم شود (مثل /profile/)، مستقیماً index.html را بده
		if strings.HasSuffix(reqPath, "/") || reqPath == "/" {
			if serveStatic(w, r, "index.html") {
				return
			}
		}

		name := strings.TrimPrefix(reqPath, "/")
		if serveStatic(w, r, name) {
			return
		}

		// 🛠️ SPA Fallback نهایی
		if serveStatic(w, r, "index.html") {
			return
		}
		http.NotFound(w, r)
	})
}
