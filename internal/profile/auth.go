package profile

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/angel/angel/internal/auth"
)

// این تابع همان چیزی است که main.go به دنبال آن می‌گردد
func RegisterAuth(mux *http.ServeMux) { mux.HandleFunc("/api/auth/token", handleToken) }

// تابع امن برای اعتبارسنجی initData تلگرام و بله
func validateInitData(initData string, botToken string) (int64, bool) {
	params, err := url.ParseQuery(initData)
	if err != nil {
		return 0, false
	}

	hash := params.Get("hash")
	if hash == "" {
		return 0, false
	}
	params.Del("hash")

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckString []string
	for _, k := range keys {
		dataCheckString = append(dataCheckString, k+"="+params.Get(k))
	}
	dataStr := strings.Join(dataCheckString, "\n")

	// ساخت کلید مخفی طبق مستندات رسمی تلگرام/بله
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	// محاسبه هش
	computedHash := hmac.New(sha256.New, secret)
	computedHash.Write([]byte(dataStr))
	computedHex := hex.EncodeToString(computedHash.Sum(nil))

	// مقایسه امن
	if !hmac.Equal([]byte(computedHex), []byte(hash)) {
		return 0, false
	}

	// استخراج آیدی کاربر
	userStr := params.Get("user")
	if userStr == "" {
		return 0, false
	}
	var tgUser struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(userStr), &tgUser); err != nil {
		return 0, false
	}

	return tgUser.ID, true
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	initData := r.URL.Query().Get("initData")
	platform := r.URL.Query().Get("platform") // "telegram" یا "bale"
	
	// برای محیط تست (در صورت نبود initData)
	if initData == "" {
		idStr := r.URL.Query().Get("id")
		if idStr != "" {
			if _, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				t, err := auth.Generate(idStr)
				if err == nil {
					writeJSON(w, 200, map[string]any{"ok": true, "token": t, "user_id": idStr})
					return
				}
			}
		}
		writeJSON(w, 400, map[string]any{"ok": false, "error": "initData required"})
		return
	}

	// انتخاب توکن بر اساس پلتفرم
	var botToken string
	if platform == "bale" {
		botToken = os.Getenv("BALE_TOKEN")
	} else {
		botToken = os.Getenv("BOT_TOKEN")
		if botToken == "" {
			botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
		}
	}

	if botToken == "" {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "bot token not configured for platform: " + platform})
		return
	}

	userID, ok := validateInitData(initData, botToken)
	if !ok {
		writeJSON(w, 401, map[string]any{"ok": false, "error": "invalid initData"})
		return
	}

	idStr := strconv.FormatInt(userID, 10)
	t, err := auth.Generate(idStr)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "token": t, "user_id": idStr})
}
