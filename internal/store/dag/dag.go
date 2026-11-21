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
	MuMap  map[string]*sync.RWMutex
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, dags: make(map[string]*workflow.DAG), MuMap: make(map[string]*sync.RWMutex)}
}

func (s *InMemoryStore) Store(dag *workflow.DAG) error {
	if mu, ok := s.MuMap[dag.WorkflowID]; ok {
		mu.Lock()
		defer mu.Unlock()
	} else {
		mu := &sync.RWMutex{}
		s.MuMap[dag.WorkflowID] = mu

		mu.Lock()
		defer mu.Unlock()
	}

	if _, ok := s.dags[dag.WorkflowID]; ok {
		return fmt.Errorf("dag %s already exists", dag.WorkflowID)
	}

	s.dags[dag.WorkflowID] = dag

	return nil
}
