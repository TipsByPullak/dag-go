package validate

import (
	"fmt"
	"log/slog"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store"
)

type Validator struct {
	logger *slog.Logger
	store  *store.Store
}

func NewValidator(logger *slog.Logger, store *store.Store) *Validator {
	return &Validator{logger: logger, store: store}
}

func (v *Validator) Idempotent(key string) error {
	// Get function uses a Read lock
	if ok := v.store.IdempotencyStore.CheckAndSet(key); !ok {
		v.logger.Warn("Idempotency store hit")
		return fmt.Errorf("idempotency key already present")
	}

	return nil
}
