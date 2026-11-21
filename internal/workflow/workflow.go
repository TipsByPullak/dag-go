package workflow

import (
	"context"
	"log/slog"
	"sync"
	"uuid"

	"github.com/TipsByPullak/phonepe-coding-round/internal/workflow/request"
)

type DAG struct {
	WorkflowID    string              `json:"workflow_id"`
	Graph         map[string][]string `json:"graph"`
	NumberOfNodes int                 `json:"number_of_nodes"`
	mu            sync.RWMutex
	State         string
	Roots         []string
}

const (
	SCHEDULED = "scheduled"
	PENDING   = "pending"
	ACCEPTED  = "accepted"
	RUNNING   = "running"
	FAILED    = "failed"
	ABORTED   = "aborted"
	SUCCEEDED = "succeeded"
)

func NewDAG(ctx context.Context, logger *slog.Logger, workflow request.CreateWorkflow) *DAG {
	dag := &DAG{
		WorkflowID:    uuid.New().String(),
		NumberOfNodes: len(workflow.Tasks),
		State:         PENDING,
	}

	dag.Graph, dag.Roots, dag.NumberOfNodes = constructGraph(workflow.Tasks)

	logger.InfoContext(ctx, "DAG Created", "workflowID", dag.WorkflowID, "graph", dag.Graph)

	return dag
}

func constructGraph(tasks []request.Task) (map[string][]string, []string, int) {
	numberOfNodes := len(tasks)
	roots := make([]string, numberOfNodes)
	graph := make(map[string][]string)

	for _, task := range tasks {
		if len(task.Dependencies) == 0 {
			roots = append(roots, task.Name)
		}

		if _, ok := graph[task.Name]; !ok {
			graph[task.Name] = make([]string, numberOfNodes)
		}
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if _, ok := graph[dep]; !ok {
				graph[dep] = make([]string, numberOfNodes)
			} else {
				graph[dep] = append(graph[dep], task.Name)
			}
		}
	}

	return graph, roots, numberOfNodes
}
