package main

type payload struct {
	Message string `json:"message"`
}

type job struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	Completed bool   `json:"completed"`
}
