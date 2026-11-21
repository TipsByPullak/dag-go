package idempotency

import (
	"log/slog"
	"sync"
)

type InMemoryStore struct {
	logger          *slog.Logger
	Mu              sync.RWMutex
	idempotencyKeys map[string]struct{}
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, idempotencyKeys: make(map[string]struct{})}
}

func (ims *InMemoryStore) CheckAndSet(key string) bool {
	ims.Mu.Lock()
	defer ims.Mu.Unlock()

	if _, ok := ims.idempotencyKeys[key]; !ok {
		ims.idempotencyKeys[key] = struct{}{}
		return true
	}

	return false
}
