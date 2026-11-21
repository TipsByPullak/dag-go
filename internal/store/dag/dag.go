package dag

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
)

type InMemoryStore struct {
	logger *slog.Logger
	dags   map[string]*workflow.DAG // stored as workflowID -> JSON
	Mu     sync.RWMutex
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, dags: make(map[string]*workflow.DAG)}
}

func (s *InMemoryStore) Store(dag *workflow.DAG) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if _, ok := s.dags[dag.WorkflowID]; ok {
		return fmt.Errorf("dag %s already exists", dag.WorkflowID)
	}

	s.dags[dag.WorkflowID] = dag

	return nil
}
