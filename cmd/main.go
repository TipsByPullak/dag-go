package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.InfoContext(ctx, "fatal", "err", err)
		return err
	}

	return nil
}
