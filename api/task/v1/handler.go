package task

import (
	middlewares "github.com/AkashGit21/task-ms/internal/middleware"
	"github.com/AkashGit21/task-ms/internal/task"
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	taskRouter := router.PathPrefix("/api/v1").Subrouter()
	taskRouter.Use(middlewares.TransactionalLogMiddleware)
	taskRouter.Use(middlewares.PanicRecoveryMiddleware)
	taskRouter.Use(middlewares.AuthMiddleware)
	task.TaskV1Handler(taskRouter)

	return router, nil
}
