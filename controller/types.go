package main

type Payload struct {
	Message string `json:"message"`
}

type StatusType string

const (
	Pending    StatusType = "pending"
	InProgress StatusType = "in progress"
	Completed  StatusType = "completed"
	Failed     StatusType = "failed"
)

type Job struct {
	Id        string     `json:"id"`
	Arguments string     `json:"-"`
	Status    StatusType `json:"status"`
	Result    string     `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
}

type CreateResponse []struct {
	NetworkInterfaces []struct {
		NetworkIP string `json:"networkIP"`
	} `json:"networkInterfaces"`
}
