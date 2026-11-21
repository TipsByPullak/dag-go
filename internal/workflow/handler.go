package workflow

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/request"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/health", h.health)
	mux.HandleFunc("POST /v1/workflows", h.ExecuteWorkflow)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "health route hit")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("OK")
}

func (h *Handler) ExecuteWorkflow(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "workflow handler hit")

	var req request.CreateWorkflow

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(r.Context(), "Workflow received", req)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("OK")
}
