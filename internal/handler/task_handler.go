package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/handrixn/task-tracker/internal/model"
	"github.com/handrixn/task-tracker/internal/service"
	"github.com/handrixn/task-tracker/internal/util"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTaskPayload model.TaskInput

	if err := json.NewDecoder(r.Body).Decode(&newTaskPayload); err != nil {
		util.JsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	validatorResult := util.ValidatePayload(newTaskPayload)
	if validatorResult != nil {
		util.JsonResponse(w, http.StatusBadRequest, validatorResult)
		return
	}

	taskCreated, err := th.taskService.Create(&newTaskPayload)

	if err != nil {
		log.Println(err)
		util.JsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	taskResponseData := model.TaskOutput{
		ID:          taskCreated.UUID,
		Title:       taskCreated.Title,
		Description: taskCreated.Description,
		DueDate:     taskCreated.DueDate.Format("2006-01-02"),
		Status:      taskCreated.Status,
	}

	util.JsonResponse(w, http.StatusCreated, taskResponseData)
}
