package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TipsByPullak/phonepe-coding-round/constants"
	"github.com/TipsByPullak/phonepe-coding-round/internal/request"
	"github.com/TipsByPullak/phonepe-coding-round/internal/worker"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/handler"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	if err := boot(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func boot(logger *slog.Logger) error {
	// main server context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// workers context
	wCtx, wStop := context.WithCancel(context.Background())
	defer wStop()

	port := constants.Port

	mux := http.NewServeMux()

	wHandler := handler.NewHandler(logger)
	wHandler.RegisterRoutes(mux)

	wp := &worker.Pool{}
	wp.Handler = wHandler

	mux.HandleFunc("POST /v1/workers/start", func(w http.ResponseWriter, r *http.Request) {
		var req request.TriggerWorkers

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		// Unmarshall request
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
			return
		}

		logger.InfoContext(r.Context(), "Trigger Workers received", "request", req)

		if req.Count == 0 {
			req.Count = constants.DefaultNumberOfWorkers // default
		}

		wp.NumWorkers = req.Count

		if err := wp.Start(wCtx, logger); err != nil {
			http.Error(w, fmt.Sprintf("failed to start workers: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("OK")
	})

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

	wStop()
	wp.Wait()

	logger.Info("workers shut down gracefully")

	return nil
}
