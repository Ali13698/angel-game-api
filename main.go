package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/games/echo"
	"github.com/angel/angel/internal/games/ludo"
	"github.com/angel/angel/internal/db"
	"github.com/angel/angel/internal/hub"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	reg := registry.New(ctx)
	if err := reg.Register(ludo.New(), echo.New()); err != nil {
		log.Fatal(err)
	}

	if err := db.Init(ctx); err != nil {
		log.Printf("[db] FAIL: %v", err)
	} else {
		defer db.Close()
	}

	h := hub.New(reg)
	br := bridge.New()
	h.SetBridge(br)

	mux := http.NewServeMux()
	mux.Handle("/bridge", br)
	mux.Handle("/ws", h)
	mux.HandleFunc("/files/", storage.ProxyHandler)
	mux.Handle("/", http.FileServer(http.Dir("/app/web")))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/dbz", func(w http.ResponseWriter, r *http.Request) {
		p := db.Pool()
		if p == nil { http.Error(w, "db not init", 500); return }
		if err := p.Ping(r.Context()); err != nil { http.Error(w, err.Error(), 500); return }
		_, _ = w.Write([]byte("db ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			GameID  string   `json:"game_id"`
			Players []string `json:"players"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		players := make([]types.UserID, len(body.Players))
		for i, p := range body.Players {
			players[i] = types.UserID(p)
		}
		rm, err := h.OpenRoom(types.GameID(body.GameID), players)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"room_id": string(rm.ID())})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "0.0.0.0:" + port

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", addr)
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal")
	case err := <-errCh:
		if err != nil {
			log.Fatalf("server: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = h.Shutdown(shutdownCtx)
	_ = reg.Close(shutdownCtx)
	log.Printf("bye")
}
