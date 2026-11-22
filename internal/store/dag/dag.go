package dag

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
)

type InMemoryStore struct {
	mu     sync.RWMutex // lock for db
	logger *slog.Logger

	Dags map[string]*workflow.DAG // stored as workflowID -> JSON
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, Dags: make(map[string]*workflow.DAG)}
}

func (s *InMemoryStore) Store(dag *workflow.DAG) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Dags[dag.WorkflowID]; exists {
		return fmt.Errorf("dag %s already exists", dag.WorkflowID)
	}

	s.Dags[dag.WorkflowID] = dag

	return nil
}

func (s *InMemoryStore) Get(workflowID string) (*workflow.DAG, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dag, ok := s.Dags[workflowID]
	return dag, ok
}

func (s *InMemoryStore) FindAvailable() *workflow.DAG {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, dag := range s.Dags {
		dag.Mu.Lock()
		if dag.State == workflow.PENDING {

			dag.Mu.Unlock()
			return dag
		}

		dag.Mu.Unlock()
	}

	return nil
}
