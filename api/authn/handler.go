package authn

import (
	"github.com/AkashGit21/task-ms/internal/task"
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	taskRouter := router.PathPrefix("/api").Subrouter()
	task.TaskV1Handler(taskRouter)

	return router, nil
}
