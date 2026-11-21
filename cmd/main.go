package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TipsByPullak/phonepe-coding-round/constants"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := boot(ctx, log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func boot(ctx context.Context, logger *slog.Logger) error {
	port := constants.Port

	mux := http.NewServeMux()

	wHandler := workflow.NewHandler(logger)
	wHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		logger.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil {
			logger.InfoContext(ctx, "fatal", "err", err)
		}
	}()

	<-ctx.Done()
	logger.Info("starting shut down gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shut down forcefully", err)
	}

	logger.Info("shut down gracefully")

	return nil
}
