package store

import (
	"log/slog"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store/dag"
	"github.com/TipsByPullak/phonepe-coding-round/internal/store/idempotency"
)

type Store struct {
	Logger           *slog.Logger
	IdempotencyStore *idempotency.InMemoryStore
	DAGStore         *dag.InMemoryStore
}

func NewStore(logger *slog.Logger) *Store {
	return &Store{
		Logger:           logger,
		IdempotencyStore: idempotency.NewInMemoryStore(logger),
		DAGStore:         dag.NewInMemoryStore(logger),
	}
}
