package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/handrixn/task-tracker/internal/handler"
	"github.com/handrixn/task-tracker/internal/middleware"
)

func NewTaskRouter(r *mux.Router, h *handler.TaskHandler) *mux.Router {
	taskRoute := r.PathPrefix("/tasks").Subrouter()
	taskRoute.Use(middleware.ValidateAPIToken)

	taskRoute.HandleFunc("/", h.TaskList).Methods(http.MethodGet)
	taskRoute.HandleFunc("/create", h.CreateTask).Methods(http.MethodPost)
	taskRoute.HandleFunc("/{uuid}/update", h.UpdateTask).Methods(http.MethodPut)
	taskRoute.HandleFunc("/summary", h.TaskSummary).Methods(http.MethodGet)

	taskRoute.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
	}).Methods(http.MethodGet)

	return taskRoute
}
