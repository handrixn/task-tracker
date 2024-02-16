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
	UpdateTask(taskID string, updateTask *model.TaskInputUpdate) (*model.Task, error)
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

func (ts *taskService) UpdateTask(taskUUID string, updateTask *model.TaskInputUpdate) (*model.Task, error) {
	task, err := ts.taskRepo.GetTaskByUUID(taskUUID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	dateTime, err := time.Parse("2006-01-02", updateTask.DueDate)

	if err != nil {
		return nil, err
	}

	task.Title = updateTask.Title
	task.Description = updateTask.Description
	task.DueDate = dateTime

	if updateTask.Status == constant.COMPLETED {
		task.Status = updateTask.Status
	}

	updatedTask, err := ts.taskRepo.UpdateTask(task.ID, task)

	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}
