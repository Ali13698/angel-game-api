package bridge

import (
	"net/http"
	"strings"
)

func isWebSocketUpgrade(r *http.Request) bool {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return false
	}
	return strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

func originAllowed(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	if origin == "null" {
		return false
	}
	host := strings.ToLower(r.Host)
	for _, suffix := range []string{"nomi98.com", "nomi98.ir"} {
		if strings.HasSuffix(host, suffix) {
			return strings.Contains(origin, suffix)
		}
	}
	return strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "http://127.0.0.1")
}
