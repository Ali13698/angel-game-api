package profile

import (
	"encoding/json"
	"net/http"

	"github.com/angel/angel/internal/db"
)

func RegisterFavorites(mux *http.ServeMux) {
	mux.HandleFunc("/api/favorites", handleGetFavorites)
	mux.HandleFunc("/api/favorites/toggle", handleToggleFavorite)
}

func handleGetFavorites(w http.ResponseWriter, r *http.Request) {
	uid, err := authFromReq(r)
	if err != nil {
		writeJSON(w, 401, map[string]any{"ok": false, "error": "unauthorized"})
		return
	}

	p := db.Pool()
	rows, err := p.Query(r.Context(), `SELECT game_id FROM user_favorites WHERE user_id = $1`, uid)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "db error"})
		return
	}
	defer rows.Close()

	var games []string
	for rows.Next() {
		var g string
		rows.Scan(&g)
		games = append(games, g)
	}
	writeJSON(w, 200, map[string]any{"ok": true, "favorites": games})
}

func handleToggleFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", 405)
		return
	}
	uid, err := authFromReq(r)
	if err != nil {
		writeJSON(w, 401, map[string]any{"ok": false, "error": "unauthorized"})
		return
	}

	var req struct {
		GameID string `json:"game_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GameID == "" {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "game_id required"})
		return
	}

	p := db.Pool()
	// بررسی وجود
	var exists bool
	p.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM user_favorites WHERE user_id=$1 AND game_id=$2)`, uid, req.GameID).Scan(&exists)

	if exists {
		_, _ = p.Exec(r.Context(), `DELETE FROM user_favorites WHERE user_id=$1 AND game_id=$2`, uid, req.GameID)
		writeJSON(w, 200, map[string]any{"ok": true, "status": "removed"})
	} else {
		_, _ = p.Exec(r.Context(), `INSERT INTO user_favorites (user_id, game_id) VALUES ($1, $2)`, uid, req.GameID)
		writeJSON(w, 200, map[string]any{"ok": true, "status": "added"})
	}
}
