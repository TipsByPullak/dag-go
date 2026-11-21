package request

type CreateWorkflow struct {
	IdempotencyKey string `json:"idempotency_key"`
	Tasks          []Task `json:"tasks"`
}

type Task struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
}

type TriggerWorkers struct {
	Count int `json:"count"`
}
