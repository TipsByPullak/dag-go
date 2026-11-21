package workflow

import (
	"context"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"uuid"

	"github.com/TipsByPullak/phonepe-coding-round/internal/request"
)

type DAG struct {
	Mu sync.RWMutex

	WorkflowID    string              `json:"workflow_id"`
	Graph         map[string][]string `json:"graph"`
	NumberOfNodes int                 `json:"number_of_nodes"`
	State         string              `json:"state"`
	Roots         []string            `json:"roots"`
}

type Task struct {
	Mu sync.RWMutex

	WorkflowID   string   `json:"workflow_id"`
	Name         string   `json:"id"`
	Dependencies []string `json:"dependencies"`
	State        string   `json:"state"`
}

const (
	SCHEDULED           = "scheduled"
	PENDING             = "pending"
	ACCEPTED            = "accepted"
	RUNNING             = "running"
	FAILED              = "failed"
	ABORTED             = "aborted"
	SUCCEEDED           = "succeeded"
	PARTIALLY_COMPLETED = "partially completed"
)

func NewDAG(ctx context.Context, logger *slog.Logger, workflow request.CreateWorkflow) *DAG {
	dag := &DAG{
		WorkflowID:    uuid.New().String(),
		NumberOfNodes: len(workflow.Tasks),
		State:         PENDING,
	}

	dag.Graph, dag.Roots, dag.NumberOfNodes = constructGraph(workflow.Tasks)

	// TODO: Change this to debug log
	logger.InfoContext(ctx, "DAG Created", "workflowID", dag.WorkflowID, "graph", dag.Graph)

	return dag
}

func constructGraph(tasks []request.Task) (map[string][]string, []string, int) {
	numberOfNodes := len(tasks)
	roots := make([]string, 0)
	graph := make(map[string][]string)

	for _, task := range tasks {
		taskName := strings.TrimSpace(task.Name)

		// Roots have no dependencies
		if len(task.Dependencies) == 0 {
			roots = append(roots, taskName)
		}

		// Fill every task name in the adjacency graph list
		if _, ok := graph[taskName]; !ok {
			graph[taskName] = make([]string, 0)
		}
	}

	// for connecting the graph
	for _, task := range tasks {
		taskName := strings.TrimSpace(task.Name)

		if len(task.Dependencies) == 0 {
			continue
		}
		for _, dep := range task.Dependencies {
			trimmedDep := strings.TrimSpace(dep)
			if !slices.Contains(graph[trimmedDep], taskName) {
				graph[trimmedDep] = append(graph[trimmedDep], taskName)
			}
		}
	}

	return graph, roots, numberOfNodes
}

func NewTask(_ context.Context, _ *slog.Logger, workflowID string, name string, deps []string) *Task {
	return &Task{
		WorkflowID:   workflowID,
		State:        PENDING,
		Name:         name,
		Dependencies: deps,
	}
}
