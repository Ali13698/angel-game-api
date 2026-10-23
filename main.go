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

	"github.com/angel/angel/internal/auth"
	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/cache"
	"github.com/angel/angel/internal/chat"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/db"
	"github.com/angel/angel/internal/events"
	"github.com/angel/angel/internal/games/echo"
	"github.com/angel/angel/internal/hub"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/profile"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/signaling"
	"github.com/angel/angel/internal/storage"
	"github.com/angel/angel/internal/wasm"
	"github.com/angel/angel/internal/web"
)

func main() {
	logger.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := wasm.Init(ctx); err != nil {
		log.Fatalf("[wasm] FAIL: %v", err)
	}
	defer wasm.Close(ctx)

	reg := registry.New(ctx)
	if err := reg.Register(echo.New()); err != nil {
		log.Fatal(err)
	}

	gamesLoaded := wasm.ListGames()
	log.Printf("[wasm] Found %d games in directory.", len(gamesLoaded))
	for _, gameID := range gamesLoaded {
		if err := reg.Register(wasm.NewHandler(gameID)); err != nil {
			log.Printf("[wasm] ❌ Failed to register %s: %v", gameID, err)
		} else {
			log.Printf("[wasm] ✅ Registered plugin: %s", gameID)
		}
	}

	if err := db.Init(ctx); err != nil {
		log.Printf("[db] FAIL: %v", err)
	} else {
		defer db.Close()
	}
	if err := cache.Init(ctx); err != nil {
		log.Printf("[cache] FAIL: %v", err)
	} else {
		defer cache.Close()
	}
	if err := storage.Init(ctx); err != nil {
		log.Printf("[storage] FAIL: %v", err)
	}
	auth.Init()

	kafkaPublisher := events.NewKafkaProducer()
	defer kafkaPublisher.Close()

	h := hub.New(reg, kafkaPublisher)
	br := bridge.New()
	h.SetBridge(br)

	chatHub := chat.NewHub()
	signalHub := signaling.NewHub()

	mux := http.NewServeMux()
	profile.Register(mux)
	profile.RegisterAuth(mux)
	
	mux.Handle("/ws", h)
	mux.Handle("/ws/chat", chatHub)
	mux.Handle("/ws/signal", signalHub)
	mux.Handle("/bridge", br)

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
			log.Printf("[rooms] ❌ Error opening room for %s: %v", body.GameID, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("[rooms] ✅ Opened room %s for game %s", rm.ID(), body.GameID)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"room_id": string(rm.ID())})
	})

	web.Register(mux)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("🚀 listening on %s", srv.Addr)
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
