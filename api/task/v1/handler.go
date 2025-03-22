package task

import (
	"github.com/AkashGit21/task-ms/internal/task"
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	taskRouter := router.PathPrefix("/api/v1").Subrouter()
	task.TaskHandler(taskRouter)

	return router, nil
}
