package store

import (
	"log/slog"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store/dag"
	"github.com/TipsByPullak/phonepe-coding-round/internal/store/idempotency"
	"github.com/TipsByPullak/phonepe-coding-round/internal/store/task"
)

type Store struct {
	Logger           *slog.Logger
	IdempotencyStore *idempotency.InMemoryStore
	DAGStore         *dag.InMemoryStore
	TaskStore        *task.InMemoryStore
}

func NewStore(logger *slog.Logger) *Store {
	return &Store{
		Logger:           logger,
		IdempotencyStore: idempotency.NewInMemoryStore(logger),
		DAGStore:         dag.NewInMemoryStore(logger),
		TaskStore:        task.NewInMemoryStore(logger),
	}
}
