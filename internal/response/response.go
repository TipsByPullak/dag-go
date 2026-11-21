package response

type CreateWorkflow struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Tasks  []Task `json:"tasks"`
}

type Task struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
