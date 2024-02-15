package service

import "github.com/handrixn/task-tracker/internal/repository"

type TaskService interface {
	// Implement service methods
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *taskService {
	return &taskService{taskRepo: taskRepo}
}
