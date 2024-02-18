package model

import "time"

type Task struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskSummary struct {
	Total      int `json:"total"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
}

type TaskInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	DueDate     string `json:"due_date" validate:"required,datetime=2006-01-02"`
}

type TaskInputUpdate struct {
	TaskInput
	Status string `json:"status" validate:"omitempty,eq=completed"`
}

type TaskOutput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}
