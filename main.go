package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/angel/angel/internal/bridge"
	"github.com/angel/angel/internal/core/types"
	"github.com/angel/angel/internal/games/echo" // فرض بر این است که بازی echo اینجاست
	"github.com/angel/angel/internal/hub"
	"github.com/angel/angel/internal/registry"
	"github.com/angel/angel/internal/wasm"
)

// echoFactory پیاده‌سازی GameFactory برای بازی echo
type echoFactory struct{}

func (f *echoFactory) ID() types.GameID { return "echo" }
func (f *echoFactory) New(ctx context.Context) (types.Game, error) {
	// اگر بازی echo استیت اشتراکی ندارد، اینجا می‌توانید مستقیم return کنید
	return echo.New(), nil
}

func main() {
	log.Println("Starting Angel Game API...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize Registry
	reg := registry.New(ctx)

	// 2. Register Games via Factories
	// ثبت بازی echo با استفاده از کارخانه جدید
	if err := reg.Register(&echoFactory{}); err != nil {
		log.Fatalf("Failed to register echo game: %v", err)
	}

	// ثبت بازی‌های WASM (لودو، دوز و...)
	// به جای NewHandler از NewFactory استفاده می‌کنیم
	gameID := "ludo-light" // یا هر بازی دیگری که در پوشه games دارید
	if err := reg.Register(wasm.NewFactory(gameID)); err != nil {
		log.Fatalf("Failed to register wasm game %s: %v", gameID, err)
	}

	// 3. Initialize Bridge & Hub
	br := bridge.New()
	h := hub.New(reg, nil) // nil در صورت عدم استفاده از Kafka در این مرحله
	h.SetBridge(br)

	// 4. Setup HTTP Server
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", h.ServeHTTP)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// 5. Graceful Shutdown setup
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("Server listening on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
