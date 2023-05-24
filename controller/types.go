package main

type Payload struct {
	Message string `json:"message"`
}

type StatusType string

const (
	Pending    StatusType = "pending"
	InProgress StatusType = "in progress"
	Completed  StatusType = "completed"
)

type Job struct {
	Id      string     `json:"id"`
	Message string     `json:"message"`
	Status  StatusType `json:"status"`
	Output  string     `json:"output,omitempty"`
}
