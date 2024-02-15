package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/handrixn/task-tracker/internal/handler"
)

func NewTaskRouter(r *mux.Router, h *handler.TaskHandler) *mux.Router {
	taskRoute := r.PathPrefix("/tasks").Subrouter()

	taskRoute.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "PONG"})
	}).Methods(http.MethodGet)
	return taskRoute
}
