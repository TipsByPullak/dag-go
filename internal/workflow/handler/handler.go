package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store"
	"github.com/TipsByPullak/phonepe-coding-round/internal/store/idempotency"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/request"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/response"
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Unmarshall request
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(r.Context(), "Workflow received", "request", req)

	// Validate request
	if err := h.Validator.ValidateRequest(&req); err != nil {
		http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create DAG
	dag := workflow.NewDAG(r.Context(), h.logger, req)

	// Check for Cycles
	if err := h.Validator.Acyclic(dag.Graph); err != nil {
		http.Error(w, "Invalid Input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check for idempotency
	if ok := h.Storage.IdempotencyStore.CheckAndSet(req.IdempotencyKey, dag.WorkflowID); !ok {
		http.Error(w, "Invalid Input: idempotency key already set", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(r.Context(), "Workflow created", "DAG", dag)

	err := h.ProcessDAG(r.Context(), dag)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := response.CreateWorkflow{
		ID:     dag.WorkflowID,
		Status: dag.State,
		Tasks:  make([]response.Task, len(dag.Graph)),
	}

	for node := range dag.Graph {
		resp.Tasks = append(resp.Tasks, response.Task{
			Name:   node,
			Status: dag.State, // TODO: Change to task.state once task is construed
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ProcessDAG(ctx context.Context, dag *workflow.DAG) error {
	err := h.Storage.DAGStore.Store(dag)
	if err != nil {
		return err
	}

	h.logger.InfoContext(ctx, "DAG Stored", "DAG", dag)
	return nil
}
