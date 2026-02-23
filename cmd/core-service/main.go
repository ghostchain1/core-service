package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ghostchain1/core-service/internal/config"
	"github.com/ghostchain1/core-service/internal/server"
	"github.com/ghostchain1/core-service/pkg/logging"
	"github.com/ghostchain1/core-service/pkg/version"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Default().Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(cfg.LogLevel)

	// Print version info
	logger.Info("starting service",
		"service", version.ServiceName,
		"version", version.Version,
		"commit", version.Commit,
		"built", version.BuildTime)

	// Create server
	srv := server.New(cfg, logger)

	// Start server in a goroutine
	go func() {
		logger.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	} else {
		logger.Info("server exited gracefully")
	}
}
