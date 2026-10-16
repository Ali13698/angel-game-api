package profile

import (
	"net/http"
	"strconv"

	"github.com/angel/angel/internal/auth"
)

func RegisterAuth(mux *http.ServeMux) { mux.HandleFunc("/api/auth/token", handleToken) }

func handleToken(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }
	if _, err := strconv.ParseInt(idStr, 10, 64); err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "bad id"}); return
	}
	t, err := auth.Generate(idStr)
	if err != nil { writeJSON(w, 500, map[string]any{"ok": false, "error": err.Error()}); return }
	writeJSON(w, 200, map[string]any{"ok": true, "token": t, "user_id": idStr})
}
