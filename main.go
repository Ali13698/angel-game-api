package main

import (
	"context"
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
	"github.com/angel/angel/internal/wire"
)

func main() {
	logger.Init()

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	reg := registry.New(ctx)
	if err := reg.Register(echo.New()); err != nil {
		log.Printf("[startup] echo registration failed: %v", err)
		return
	}

	if err := wasm.Init(ctx); err != nil {
		log.Printf("[startup] WASM initialization failed: %v", err)
		wasm.Close(context.Background())
		return
	}
	defer wasm.Close(context.Background())

	registered := 0
	for _, gameID := range wasm.ListGames() {
		if err := reg.Register(wasm.NewHandler(gameID)); err != nil {
			log.Printf("[wasm] registration failed %s: %v", gameID, err)
			continue
		}
		registered++
		log.Printf("[wasm] registered: %s", gameID)
	}
	log.Printf("[wasm] registered games: %d", registered)

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
	h.StatsRecorder = wire.NewStatsAdapter()

	br := bridge.New()
	h.SetBridge(br)

	mux := http.NewServeMux()
	profile.Register(mux)
	profile.RegisterAuth(mux)
	profile.RegisterFavorites(mux)

	mux.Handle("/ws", h)
	mux.Handle("/ws/chat", chat.NewHub())
	mux.Handle("/ws/signal", signaling.NewHub())
	mux.Handle("/bridge", br)

	registerRooms(mux, h)
	web.Register(mux)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if ctx.Err() != nil || registered == 0 {
			http.Error(w, "games unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ready\n"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              "0.0.0.0:" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("[http] listening on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Printf("[http] shutdown requested")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[http] server failed: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx)
	_ = h.Shutdown(shutdownCtx)
	_ = reg.Close(shutdownCtx)
	log.Printf("[http] stopped")
}
