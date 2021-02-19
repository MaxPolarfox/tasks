package types

type Task struct {
	ID string `json:"id"`
	Data string `json:"data"`
}

type CreatedTaskResponse struct {
	ID string `json:"id"`
}