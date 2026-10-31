package main

import (
	"net/http"

	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/chat"
	"github.com/angel/angel/internal/hub"
	"github.com/angel/angel/internal/profile"
	"github.com/angel/angel/internal/signaling"
	"github.com/angel/angel/internal/web"
)

func newMux(h *hub.Hub, br *bridge.Bridge, ready func() bool) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/ws", h)
	mux.Handle("/bridge", br)
	mux.Handle("/ws/chat", chat.NewHub())
	mux.Handle("/ws/signal", signaling.NewHub())

	registerRooms(mux, h)
	profile.Register(mux)
	profile.RegisterAuth(mux)
	profile.RegisterFavorites(mux)
	web.Register(mux)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if ready != nil && !ready() {
			http.Error(w, "games unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ready\n"))
	})

	return mux
}
