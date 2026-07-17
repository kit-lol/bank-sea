package main

import (
	"bank-sea/internal/pkg/config"
	"bank-sea/internal/pkg/logger"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Configuration initialization
	cfg := config.MustLoad()

	// 2. Logger initialization
	log := logger.SetupLogger(cfg.Env)
	log.Info("Bank server initialization", slog.String("env", cfg.Env))

	// 3. Router setup
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Received request for /health")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 4. HTTP server setup
	srv := &http.Server{
		Addr:         ":" + cfg.HttpServer.Port,
		Handler:      mux,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	// 5. Starting the server in a separate thread (goroutine)
	go func() {
		log.Info("The server is running", slog.String("port", cfg.HttpServer.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 6. Waiting for a stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Stop signal received; initiating shutdown...")

	// 7. Server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Error stoping the server", slog.String("error", err.Error()))
	}

	log.Info("The server has been successfully stopped")
}
