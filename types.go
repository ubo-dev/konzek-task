package main

import "math/rand"

type Status int

const (
	OPEN Status = iota
	CLOSED
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func NewAccount(title, description string, status Status) *Task {
	return &Task{
		ID:          rand.Intn(10000),
		Title:       title,
		Description: description,
		Status:      0,
	}
}
