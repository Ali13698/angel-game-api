package profile

import (
	"context"
	"log"
	"net/http"

	"github.com/angel/angel/internal/db"
)

type deviceReq struct {
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

	uid, err := authFromReq(r)
	if err != nil {
		writeJSON(w, 401, map[string]any{"ok": false, "error": "unauthorized"})
		return
	}

	var req deviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Fingerprint == "" {
		writeJSON(w, 400, map[string]any{"ok": false, "error": "fingerprint required"})
		return
	}

	p := db.Pool()
	if p == nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "db error"})
		return
	}

	ctx := r.Context()

	// ۱. چک می‌کنیم آیا این اثر انگشت قبلاً برای کاربر دیگری ثبت شده است؟
	var existingID int64
	err = p.QueryRow(ctx, `SELECT id FROM users WHERE device_fingerprint = $1 AND id != $2 LIMIT 1`, req.Fingerprint, uid).Scan(&existingID)

	if err == nil {
		// 🚨 بات/مولتی اکانت تشخیص داده شد! (اثر انگشت تکراری)
		log.Printf("[ANTI-BOT] Multi-account detected for user %d (duplicate of %d)", uid, existingID)
		
		// افزایش ریسک و ثبت فلگ
		_, _ = p.Exec(ctx, `UPDATE users SET risk_score = risk_score + 50, risk_flags = risk_flags || ',multi_account' WHERE id = $1`, uid)
		writeJSON(w, 403, map[string]any{"ok": false, "error": "device already registered to another account"})
		return
	}

	// ۲. اگر تکراری نبود، اثر انگشت را برای این کاربر ثبت می‌کنیم
	_, err = p.Exec(ctx, `UPDATE users SET device_fingerprint = $1 WHERE id = $2`, req.Fingerprint, uid)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "failed to save fingerprint"})
		return
	}

	writeJSON(w, 200, map[string]any{"ok": true, "message": "device registered"})
}
