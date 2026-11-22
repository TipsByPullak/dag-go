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
	mu     sync.RWMutex
}

func NewInMemoryStore(logger *slog.Logger) *InMemoryStore {
	return &InMemoryStore{logger: logger, tasks: make(map[string]*workflow.Task)}
}

func (s *InMemoryStore) Store(task *workflow.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	uid := task.WorkflowID + "|" + task.Name

	if _, ok := s.tasks[uid]; ok {
		return fmt.Errorf("task %s already exists for workflow %s", task.Name, task.WorkflowID)
	}

	s.tasks[uid] = task

	return nil
}

func (s *InMemoryStore) Get(workflowID string, taskID string) (*workflow.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[workflowID+"|"+taskID]
	return task, ok
}
