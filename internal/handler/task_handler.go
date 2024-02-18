package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/handrixn/task-tracker/internal/constant"
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

func (th *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["uuid"]

	var updateTask model.TaskInputUpdate

	if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
		util.JsonResponse(w, http.StatusBadRequest, nil)
		return
	}

	validatorResult := util.ValidatePayload(updateTask)
	if validatorResult != nil {
		util.JsonResponse(w, http.StatusBadRequest, validatorResult)
		return
	}

	updatedTask, err := th.taskService.UpdateTask(taskID, &updateTask)
	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	taskResponseData := model.TaskOutput{
		ID:          updatedTask.UUID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		DueDate:     updatedTask.DueDate.Format("2006-01-02"),
		Status:      updatedTask.Status,
	}

	util.JsonResponse(w, http.StatusOK, taskResponseData)
}

func (th *TaskHandler) TaskList(w http.ResponseWriter, r *http.Request) {
	searchFilter := r.URL.Query().Get("search")
	statusFilter := r.URL.Query().Get("status")
	dueDateFilter := r.URL.Query().Get("due_date")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr == "0" || pageStr == "" {
		pageStr = "1"
	}

	if limitStr == "0" || limitStr == "" {
		limitStr = "10"
	}

	params := map[string]string{
		constant.TASK_FILTER_TITILE_NAME:   searchFilter,
		constant.TASK_FILTER_STATUS_NAME:   statusFilter,
		constant.TASK_FILTER_DUE_DATE_NAME: dueDateFilter,
		constant.TASK_LIST_PAGE_NAME:       pageStr,
		constant.TASK_LIST_LIMIT_NAME:      limitStr,
	}

	tasks, totalData, err := th.taskService.ListTasks(params)

	if err != nil {
		util.JsonResponse(w, http.StatusInternalServerError, nil)
		return
	}

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	page, _ := strconv.ParseInt(pageStr, 10, 64)

	totalPages := totalData / limit
	if totalData%limit != 0 {
		totalPages++
	}

	taskResponseData := []model.TaskOutput{}

	for _, task := range tasks {
		to := model.TaskOutput{
			ID:          task.UUID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate.Format("2006-01-02"),
			Status:      task.Status,
		}

		taskResponseData = append(taskResponseData, to)
	}

	responseData := struct {
		Tasks       []model.TaskOutput `json:"tasks"`
		TotalCount  int64              `json:"total_count"`
		TotalPages  int64              `json:"total_pages"`
		CurrentPage int64              `json:"current_page"`
	}{
		Tasks:       taskResponseData,
		TotalCount:  totalData,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	util.JsonResponse(w, http.StatusOK, responseData)
}
