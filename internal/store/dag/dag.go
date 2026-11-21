package dag

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
)

type InMemoryStore struct {
	logger *slog.Logger
	Dags   map[string]*workflow.DAG // stored as workflowID -> JSON
	MuMap  map[string]*sync.RWMutex
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, Dags: make(map[string]*workflow.DAG), MuMap: make(map[string]*sync.RWMutex)}
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

	if _, ok := s.Dags[dag.WorkflowID]; ok {
		return fmt.Errorf("dag %s already exists", dag.WorkflowID)
	}

	s.Dags[dag.WorkflowID] = dag

	return nil
}

func (s *InMemoryStore) FindAvailable() *workflow.DAG {
	for _, dag := range s.Dags {
		s.logger.Info("dag fetched", dag)
		if dag.State == workflow.PENDING {
			dag.Mu.Lock()

			// Refresh after mutex acquire
			if dag.State == workflow.PENDING {
				dag.State = workflow.RUNNING

				return dag
			}

			dag.Mu.Unlock()
		}
	}

	return nil
}
