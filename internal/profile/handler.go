package profile

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/angel/angel/internal/auth"
	"github.com/angel/angel/internal/db"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/profile/", handleProfile)
	mux.HandleFunc("/api/me", handleMe)
	mux.HandleFunc("/api/me/upload", handleUpload)
	mux.HandleFunc("/api/profile/like", handleLike) 
	log.Printf("[PROFILE] Routes registered")
}

func handleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	// خطای declared and not used اینجا بود. با _ جایگزین شد
	_, err := authFromReq(r)
	if err != nil {
		writeJSON(w, 401, map[string]any{"ok": false, "error": "unauthorized"})
		return
	}

	var req struct {
		TargetID int64 `json:"target_id"`
		IsLiked  bool  `json:"is_liked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "bad json"})
		return
	}

	p := db.Pool()
	if p == nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "db error"})
		return
	}

	diff := 1
	if !req.IsLiked {
		diff = -1
	}

	_, err = p.Exec(r.Context(), `UPDATE users SET likes_count = GREATEST(0, likes_count + $1) WHERE id = $2`, diff, req.TargetID)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "failed to update like"})
		return
	}

	writeJSON(w, 200, map[string]any{"ok": true})
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

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "GET only", 405); return }
	param := strings.TrimPrefix(r.URL.Path, "/api/profile/")
	if param == "" || param == "like" { return } 
	if param == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }

	var u *User
	var err error

	id, parseErr := strconv.ParseInt(param, 10, 64)
	if parseErr == nil {
		u, err = GetByID(r.Context(), id)
	} else {
		u, err = GetByPublicToken(r.Context(), param)
		if err != nil {
			u, err = GetByGameUsername(r.Context(), param)
		}
	}

	if err != nil { 
		writeJSON(w, 404, map[string]any{"ok": false, "error": "user not found"}); 
		return 
	}
	
	stats, _ := GetGameStats(r.Context(), u.ID)
	writeJSON(w, 200, map[string]any{"ok": true, "user": u, "stats": stats})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	uid, err := authFromReq(r)
	if err != nil { writeJSON(w, 401, map[string]any{"ok": false, "error": err.Error()}); return }

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
			writeJSON(w, 400, map[string]any{"ok": false, "error": "bad json"}); return
		}
		if err := UpdateMe(r.Context(), uid, &req); err != nil {
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
func errStr(s string) error    { return strErr(s) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
