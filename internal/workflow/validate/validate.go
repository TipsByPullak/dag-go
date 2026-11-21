package validate

import (
	"fmt"
	"log/slog"

	"github.com/TipsByPullak/phonepe-coding-round/internal/store"
	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/request"
)

type Validator struct {
	logger *slog.Logger
	store  *store.Store
}

func NewValidator(logger *slog.Logger, store *store.Store) *Validator {
	return &Validator{logger: logger, store: store}
}

func (v *Validator) Idempotent(key string, val string) error {
	if key == "" {
		v.logger.Info("Idempotency key is empty")
		return nil
	}

	// CheckAndSet function uses a lock
	if ok := v.store.IdempotencyStore.CheckAndSet(key, val); !ok {
		v.logger.Warn("Idempotency store hit")
		return fmt.Errorf(
			"idempotency key already present with workflow ID: %s",
			v.store.IdempotencyStore.Get(key),
		)
	}

	return nil
}

func (v *Validator) Acyclic(graph map[string][]string) error {
	// 0 - not visited, 1 - visiting, 2 - visited
	visited := make(map[string]int, len(graph))

	// DFS function
	var visit func(node string) bool
	visit = func(node string) bool {
		visited[node] = 1 // mark visiting

		for _, child := range graph[node] {
			if visited[child] == 1 {
				// cycle detected
				return true
			} else if visited[child] == 0 {
				if visit(child) {
					return true
				}
			}
		}

		return false
	}

	for node := range graph {
		if ok := visited[node]; ok == 0 {
			if visit(node) {
				v.logger.Warn("Cycle in DAG detected")
				return fmt.Errorf("cycle detected in graph")
			}
		}
	}

	return nil
}

func (v *Validator) ValidateRequest(req *request.CreateWorkflow) error {

	var duplicates map[string]int

	for _, task := range req.Tasks {
		if task.Name == "" {
			return fmt.Errorf("task name is empty")
		}

		if _, ok := duplicates[task.Name]; ok {
			return fmt.Errorf("duplicate task name: %s", task.Name)
		}
		duplicates[task.Name] = 1

		var dependentDuplicates map[string]int
		for _, child := range task.Dependencies {
			if child == "" {
				return fmt.Errorf("dependency name is empty in task: %s", task.Name)
			}

			if _, ok := dependentDuplicates[child]; ok {
				return fmt.Errorf("dependency name %s duplicated in task %s", child, task.Name)
			}
			dependentDuplicates[child] = 1
		}
	}

	return nil
}
