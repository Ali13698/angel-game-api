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
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { http.Error(w, "GET only", 405); return }
	param := strings.TrimPrefix(r.URL.Path, "/api/profile/")
	if param == "" { writeJSON(w, 400, map[string]any{"ok": false, "error": "id required"}); return }

	var u *User
	var err error

	// اگر عدد بود -> ID، اگر نبود -> توکن اختصاصی یا نام کاربری
	id, parseErr := strconv.ParseInt(param, 10, 64)
	if parseErr == nil {
		u, err = GetByID(r.Context(), id)
	} else {
		// جستجو با توکن اختصاصی (اولویت بالاتر از نام کاربری)
		u, err = GetByPublicToken(r.Context(), param)
		if err != nil {
			// اگر با توکن پیدا نشد، با نام کاربری جستجو کن
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

// ... (بقیه توابع handleMe, authFromReq, writeJSON و updateReq بدون تغییر)
