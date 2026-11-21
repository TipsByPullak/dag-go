package workflow

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store"
	"github.com/TipsByPullak/phonepe-coding-round/internal/store/idempotency"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/request"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/validate"
)

type Handler struct {
	logger           *slog.Logger
	Storage          *store.Store
	IdempotencyStore *idempotency.InMemoryStore
	Validator        *validate.Validator
}

func NewHandler(logger *slog.Logger) *Handler {
	newStorage := store.NewStore(logger)

	return &Handler{
		logger:    logger,
		Storage:   newStorage,
		Validator: validate.NewValidator(logger, newStorage),
	}
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

	h.logger.InfoContext(r.Context(), "Workflow received", "request", req)

	dag := NewDAG(r.Context(), h.logger, req)
	// -------------- VALIDATIONS --------------
	if err := h.Validator.Idempotent(req.IdempotencyKey); err != nil {
		http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Add validation for DAG
	// -------------- VALIDATIONS --------------

	h.logger.InfoContext(r.Context(), "Workflow created", "DAG", dag)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("OK")
}
