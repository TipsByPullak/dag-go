package idempotency

import (
	"log/slog"
	"sync"
)

type InMemoryStore struct {
	logger          *slog.Logger
	Mu              sync.RWMutex
	idempotencyKeys map[string]string
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, idempotencyKeys: make(map[string]string)}
}

func (ims *InMemoryStore) CheckAndSet(key string, val string) bool {
	ims.Mu.Lock()
	defer ims.Mu.Unlock()

	if _, ok := ims.idempotencyKeys[key]; !ok {
		ims.idempotencyKeys[key] = val
		return true
	}

	return false
}

func (ims *InMemoryStore) Get(key string) string {
	ims.Mu.RLock()
	defer ims.Mu.RUnlock()

	return ims.idempotencyKeys[key]
}
