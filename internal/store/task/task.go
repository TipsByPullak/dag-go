package task

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow"
)

type InMemoryStore struct {
	logger *slog.Logger
	tasks  map[string]*workflow.Task // stored as workflowID -> JSON
	MuMap  map[string]*sync.RWMutex
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, tasks: make(map[string]*workflow.Task), MuMap: make(map[string]*sync.RWMutex)}
}

func (s *InMemoryStore) Store(task *workflow.Task) error {
	uid := task.WorkflowID + task.Name

	if mu, ok := s.MuMap[uid]; ok {
		mu.Lock()
		defer mu.Unlock()
	} else {
		mu := &sync.RWMutex{}
		s.MuMap[uid] = mu

		mu.Lock()
		defer mu.Unlock()
	}

	if _, ok := s.tasks[uid]; ok {
		return fmt.Errorf("task %s already exists for workflow %s", task.Name, task.WorkflowID)
	}

	s.tasks[uid] = task

	return nil
}
