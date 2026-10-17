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
	log.Printf("[API] GET %s", r.URL.Path)
	if r.Method != http.MethodGet { http.Error(w, "GET only", 405); return }
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/profile/"), "/")
	if len(parts) == 0 || parts[0] == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil { writeJSON(w, 400, map[string]any{"ok": false, "error": "bad id"}); return }
	u, err := GetByID(r.Context(), id)
	if err != nil { writeJSON(w, 404, map[string]any{"ok": false, "error": err.Error()}); return }
	
	stats, _ := GetGameStats(r.Context(), id)
	writeJSON(w, 200, map[string]any{"ok": true, "user": u, "stats": stats})
}

type updateReq struct {
	DisplayName  *string `json:"display_name,omitempty"`
	GameUsername *string `json:"game_username,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
	NameEmoji    *string `json:"name_emoji,omitempty"`
	ProfileColor *string `json:"profile_color,omitempty"`
	BannerID     *string `json:"banner_id,omitempty"`
	FrameID      *string `json:"frame_id,omitempty"`
	HatID        *string `json:"hat_id,omitempty"`
	ChatBubbleID *string `json:"chat_bubble_id,omitempty"`
	BadgeID      *string `json:"badge_id,omitempty"`
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	log.Printf("[API] %s %s", r.Method, r.URL.Path)
	uid, err := authFromReq(r)
	if err != nil { 
		log.Printf("[API] Auth failed: %v", err)
		writeJSON(w, 401, map[string]any{"ok": false, "error": err.Error()}); 
		return 
	}

	if r.Method == http.MethodGet {
		u, err := GetByID(r.Context(), uid)
		if err != nil { writeJSON(w, 404, map[string]any{"ok": false, "error": err.Error()}); return }
		stats, _ := GetGameStats(r.Context(), uid)
		writeJSON(w, 200, map[string]any{"ok": true, "user": u, "stats": stats})
		return
	}
	if r.Method == http.MethodPatch {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
		var req updateReq
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("[API] Bad JSON: %v", err)
			writeJSON(w, 400, map[string]any{"ok": false, "error": "bad json"}); return
		}
		if err := UpdateMe(r.Context(), uid, &req); err != nil {
			log.Printf("[API] UpdateMe failed: %v", err)
			writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()}); return
		}
		u, _ := GetByID(r.Context(), uid)
		stats, _ := GetGameStats(r.Context(), uid)
		writeJSON(w, 200, map[string]any{"ok": true, "user": u, "stats": stats})
		return
	}
	http.Error(w, "method not allowed", 405)
}

func authFromReq(r *http.Request) (int64, error) {
	h := r.Header.Get("Authorization")
	if h == "" { return 0, errStr("missing auth") }
	token := strings.TrimPrefix(h, "Bearer ")
	claims, err := auth.Verify(token)
	if err != nil { return 0, errStr("invalid token") }
	return strconv.ParseInt(claims.UserID, 10, 64)
}

type strErr string
func (e strErr) Error() string { return string(e) }
func errStr(s string) error { return strErr(s) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}