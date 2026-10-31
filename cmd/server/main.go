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
	"github.com/angel/angel/internal/db"
	"github.com/angel/angel/internal/events"
	"github.com/angel/angel/internal/hub"
	"github.com/angel/angel/internal/logger"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/storage"
	"github.com/angel/angel/internal/wasm"
	"github.com/angel/angel/internal/wire"
)

func main() {
	logger.Init()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	reg := registry.New(ctx)

	games, err := registerGames(ctx, reg)
	if err != nil {
		log.Fatalf("[startup] %v", err)
	}
	defer wasm.Close(context.Background())
	if games == 0 {
		log.Printf("[startup] WARN: no games registered — /readyz returns 503")
	}

	if err := storage.Init(ctx); err != nil {
		log.Fatalf("[storage] %v", err)
	}

	if err := db.Init(ctx); err != nil {
		log.Printf("[db] DEGRADED: %v", err)
	} else {
		defer db.Close()
	}
	if err := cache.Init(ctx); err != nil {
		log.Printf("[cache] DEGRADED: %v", err)
	} else {
		defer cache.Close()
	}

	auth.Init()

	publisher := events.NewKafkaProducer()
	defer publisher.Close()

	h := hub.New(reg, publisher)
	// h.StatsRecorder = wire.NewStatsAdapter() // ⚠️ اگر این خط ارور داد، کامنتش کن

	br := bridge.New()
	h.SetBridge(br)

	ready := func() bool { return ctx.Err() == nil && games > 0 }
	mux := newMux(h, br, ready)

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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = h.Shutdown(shutdownCtx)
	_ = reg.Close(shutdownCtx)
	log.Printf("[http] stopped")
}
