package workflow

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/health", h.health)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "health route hit")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("OK")
}
