package profile

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/angel/angel/internal/auth"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/profile/", handleProfile)
	mux.HandleFunc("/api/me", handleMe)
	mux.HandleFunc("/api/me/upload", handleUpload)
	log.Printf("[PROFILE] Routes registered")
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "GET only", 405); return }
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/profile/"), "/")
	if len(parts) == 0 || parts[0] == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }
	
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil { writeJSON(w, 400, map[string]any{"ok": false, "error": "bad id format"}); return }
	
	log.Printf("[PROFILE] API Request for user ID: %d", id)

	u, err := GetByID(r.Context(), id)
	if err != nil { 
		log.Printf("[PROFILE] DB Error for ID %d: %v", id, err)
		writeJSON(w, 404, map[string]any{"ok": false, "error": err.Error()}); 
		return 
	}
	
	stats, err := GetGameStats(r.Context(), id)
	if err != nil {
		log.Printf("[PROFILE] Stats Error for ID %d: %v", id, err)
	}
	
	log.Printf("[PROFILE] Successfully fetched user ID: %d", id)
	writeJSON(w, 200, map[string]any{"ok": true, "user": u, "stats": stats})
}

// ... بقیه کدها (handleMe, handleUpload, authFromReq, etc) بدون تغییر ...
