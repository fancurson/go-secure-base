// Package app initializes and starts the application server.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fancurson/go-secure-base/internal/config"
	"github.com/fancurson/go-secure-base/internal/httpsrv"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

// Run initializes and starts the application server with graceful shutdown support.
func Run() error {
	cfg := config.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	serverErrors := make(chan error, 1)

	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	defer stopOngoingGracefully()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if isShuttingDown.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	})

	srv := httpsrv.NewServer(ongoingCtx, cfg, mux)

	go func() {
		logger.Info("server starting", slog.String("addr", cfg.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("listen and serve: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("critical server error: %w", err)

	case <-rootCtx.Done():
		stop()
		isShuttingDown.Store(true)
		logger.Info("Received shutdown signal, shutting down.")

		// Give time for readiness check to propagate
		time.Sleep(_readinessDrainDelay)
		logger.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		stopOngoingGracefully()
		if err != nil {
			logger.Error("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
			time.Sleep(_shutdownHardPeriod)
		}

		logger.Info("Server shut down gracefully.")
	}

	return nil
}

// Readiness endpoint
// http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
// 	if isShuttingDown.Load() {
// 		http.Error(w, "Shutting down", http.StatusServiceUnavailable)
// 		return
// 	}
// 	fmt.Fprintln(w, "OK")
// })
