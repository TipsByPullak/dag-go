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
	return &InMemoryStore{logger: logger}
}

func (ims *InMemoryStore) Set(key string) {
	ims.Mu.Lock()
	defer ims.Mu.Unlock()

	ims.idempotencyKeys[key] = struct{}{}
}

func (ims *InMemoryStore) Del(key string) {
	ims.Mu.Lock()
	defer ims.Mu.Unlock()

	delete(ims.idempotencyKeys, key)
}

func (ims *InMemoryStore) Keys() []string {
	ims.Mu.RLock()
	defer ims.Mu.RUnlock()

	keys := make([]string, 0, len(ims.idempotencyKeys))
	for k := range ims.idempotencyKeys {
		keys = append(keys, k)
	}
	return keys
}

func (ims *InMemoryStore) Len() int {
	ims.Mu.RLock()
	defer ims.Mu.RUnlock()
	return len(ims.idempotencyKeys)
}

func (ims *InMemoryStore) Get(key string) interface{} {
	ims.Mu.RLock()
	defer ims.Mu.RUnlock()
	return ims.idempotencyKeys[key]
}
