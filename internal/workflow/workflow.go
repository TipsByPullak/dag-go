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

	// TODO: Change this to debug log
	logger.InfoContext(ctx, "DAG Created", "workflowID", dag.WorkflowID, "graph", dag.Graph)

	return dag
}

func constructGraph(tasks []request.Task) (map[string][]string, []string, int) {
	numberOfNodes := len(tasks)
	roots := make([]string, 0)
	graph := make(map[string][]string)

	for _, task := range tasks {
		// Roots have no dependencies
		if len(task.Dependencies) == 0 {
			roots = append(roots, task.Name)
		}

		// Fill every task name in the adjacency graph list
		if _, ok := graph[task.Name]; !ok {
			graph[task.Name] = make([]string, 0)
		}
	}

	// for connecting the graph
	for _, task := range tasks {
		if len(task.Dependencies) == 0 {
			continue
		}
		for _, dep := range task.Dependencies {
			graph[dep] = append(graph[dep], task.Name)
		}
	}

	return graph, roots, numberOfNodes
}
