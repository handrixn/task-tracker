package service

import (
	"log"
	"time"

	"github.com/handrixn/task-tracker/internal/constant"
	"github.com/handrixn/task-tracker/internal/model"
	"github.com/handrixn/task-tracker/internal/repository"
)

type TaskService interface {
	Create(*model.TaskInput) (*model.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *taskService {
	return &taskService{taskRepo: taskRepo}
}

func (ts *taskService) Create(ti *model.TaskInput) (*model.Task, error) {
	dateTime, err := time.Parse("2006-01-02", ti.DueDate)
	if err != nil {
		log.Println("Error parsing date:", err)
		return nil, err
	}

	task := &model.Task{
		Title:       ti.Title,
		Description: ti.Description,
		DueDate:     dateTime,
		Status:      constant.IN_PROGRESS,
	}

	result, err := ts.taskRepo.Create(task)

	if err != nil {
		return nil, err
	}

	return result, nil
}
