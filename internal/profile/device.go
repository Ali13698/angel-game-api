package profile

import (
	"encoding/json" // ⬅️ این import جا افتاده بود
	"log"
	"net/http"

	"github.com/angel/angel/internal/db"
)

type deviceReq struct {
	UserID      int64  `json:"user_id"`
	Fingerprint string `json:"fingerprint"`
}

func RegisterDevice(mux *http.ServeMux) {
	mux.HandleFunc("/api/device/register", handleDeviceRegister)
}

func handleDeviceRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req deviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Fingerprint == "" || req.UserID == 0 {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "user_id and fingerprint required"})
		return
	}

	p := db.Pool()
	if p == nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "db error"})
		return
	}

	ctx := r.Context()

	var existingID int64
	err := p.QueryRow(ctx, `SELECT id FROM users WHERE device_fingerprint = $1 AND id != $2 LIMIT 1`, req.Fingerprint, req.UserID).Scan(&existingID)

	if err == nil {
		log.Printf("[ANTI-BOT] Device shared! User %d is using device previously used by %d", req.UserID, existingID)
		_, _ = p.Exec(ctx, `UPDATE users SET risk_score = risk_score + 20, risk_flags = COALESCE(risk_flags,'') || ',shared_device' WHERE id = $1`, req.UserID)
		writeJSON(w, 200, map[string]any{"ok": true, "message": "device registered with warning", "warning": "shared_device"})
		return
	}

	_, err = p.Exec(ctx, `UPDATE users SET device_fingerprint = $1 WHERE id = $2`, req.Fingerprint, req.UserID)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "failed to save fingerprint"})
		return
	}

	writeJSON(w, 200, map[string]any{"ok": true, "message": "device registered"})
}
