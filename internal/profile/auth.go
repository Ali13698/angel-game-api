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
	"time"

	"github.com/angel/angel/internal/auth"
)

func RegisterAuth(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/token", handleToken)
}

func validateInitData(initData, botToken string) (int64, bool) {
	if botToken == "" || len(initData) == 0 || len(initData) > 16384 {
		return 0, false
	}

	params, err := url.ParseQuery(initData)
	if err != nil {
		return 0, false
	}

	for _, values := range params {
		if len(values) != 1 {
			return 0, false
		}
	}

	expected, err := hex.DecodeString(params.Get("hash"))
	if err != nil || len(expected) != sha256.Size {
		return 0, false
	}

	authDate, err := strconv.ParseInt(params.Get("auth_date"), 10, 64)
	if err != nil {
		return 0, false
	}

	now := time.Now().Unix()
	if authDate < now-300 || authDate > now+30 {
		return 0, false
	}

	params.Del("hash")
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+params.Get(key))
	}

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	_, _ = secret.Write([]byte(botToken))

	signature := hmac.New(sha256.New, secret.Sum(nil))
	_, _ = signature.Write([]byte(strings.Join(lines, "\n")))

	if !hmac.Equal(signature.Sum(nil), expected) {
		return 0, false
	}

	var user struct {
		ID int64 `json:"id"`
	}
	if json.Unmarshal([]byte(params.Get("user")), &user) != nil || user.ID <= 0 {
		return 0, false
	}

	return user.ID, true
}

// GameIdentity returns only a verified identity, never a query-string user ID.
func GameIdentity(initData, platform string) (string, bool) {
	var token string

	switch platform {
	case "", "telegram":
		token = os.Getenv("BOT_TOKEN")
		if token == "" {
			token = os.Getenv("TELEGRAM_BOT_TOKEN")
		}
	case "bale":
		token = os.Getenv("BALE_TOKEN")
	default:
		return "", false
	}

	id, ok := validateInitData(initData, token)
	if !ok {
		return "", false
	}
	return strconv.FormatInt(id, 10), true
}

func handleToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	var initData, platform string

	switch r.Method {
	case http.MethodGet:
		initData = r.URL.Query().Get("initData")
		platform = r.URL.Query().Get("platform")
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 32768)
		var body struct {
			InitData string `json:"initData"`
			Platform string `json:"platform"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			writeJSON(w, 400, map[string]any{"ok": false, "error": "invalid request"})
			return
		}
		initData, platform = body.InitData, body.Platform
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := GameIdentity(initData, platform)
	if !ok {
		writeJSON(w, 401, map[string]any{
			"ok": false, "error": "invalid or expired initData",
		})
		return
	}

	token, err := auth.Generate(id)
	if err != nil {
		writeJSON(w, 500, map[string]any{"ok": false, "error": "token unavailable"})
		return
	}

	writeJSON(w, 200, map[string]any{
		"ok": true, "token": token, "user_id": id,
	})
}
