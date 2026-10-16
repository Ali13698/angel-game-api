package profile

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func Register(mux *http.ServeMux) { mux.HandleFunc("/api/profile/", handleProfile) }

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "GET only", 405); return }
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/profile/"), "/")
	if len(parts) == 0 || parts[0] == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil { writeJSON(w, 400, map[string]any{"ok": false, "error": "bad id"}); return }
	u, err := GetByID(r.Context(), id)
	if err != nil { writeJSON(w, 404, map[string]any{"ok": false, "error": err.Error()}); return }
	writeJSON(w, 200, map[string]any{"ok": true, "user": u})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
