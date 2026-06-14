package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zynu/server/internal/handlers"
	"github.com/zynu/server/internal/middleware"
	"github.com/zynu/server/pkg/config"
	"github.com/zynu/server/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok","service":"zynu-server"}`)
	})

	// Video routes
	vh := handlers.NewVideoHandler(cfg, log)
	mux.Handle("/stream/", middleware.Auth(cfg.APISecret)(http.StripPrefix("/stream", vh.StreamRouter())))
	mux.Handle("/upload/", middleware.Auth(cfg.APISecret)(http.StripPrefix("/upload", vh.UploadRouter())))

	// Webhook
	wh := handlers.NewWebhookHandler(cfg, log)
	mux.Handle("/webhook/", middleware.RateLimit(20)(wh.Router()))

	handler := middleware.CORS(middleware.Logger(log)(mux))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Infof("ZynU server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Info("Server stopped")
}
