package task

import (
	"github.com/AkashGit21/task-ms/lib/persistence/mysql"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	mysql.TaskOps
}

func NewTaskAPIHandler() *TaskHandler {
	persistenceDB, err := mysql.NewTaskPersistenceLayer()
	if err != nil {
		panic(err)
	}

	return &TaskHandler{
		persistenceDB,
	}
}

func TaskV1Handler(r *mux.Router) {
	th := NewTaskAPIHandler()

	r.HandleFunc("/tasks", th.createTask).Methods("POST")
}
