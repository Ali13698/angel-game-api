package storage

import (
	"io"
	"net/http"
	"os"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	target := os.Getenv("MINIO_INTERNAL_URL")
	if target == "" {
		target = "http://angel-cache-yti-service:9000"
	}

	// path format: /files/<bucket>/<object>
	path := r.URL.Path
	upstream := target + path

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, upstream, nil)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
