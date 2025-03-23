package task

import (
	"net/http"

	"github.com/AkashGit21/task-ms/lib/persistence/mysql"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	mysql.TaskOps
}

func newTaskAPIHandler() *TaskHandler {
	persistenceDB, err := mysql.NewTaskPersistenceLayer()
	if err != nil {
		panic(err)
	}

	return &TaskHandler{
		persistenceDB,
	}
}

func TaskV1Handler(r *mux.Router) {
	th := newTaskAPIHandler()

	r.HandleFunc("/tasks", th.createTask).Methods("POST")
	r.HandleFunc("/tasks/{taskID}", th.getTask).Methods("GET")
	r.HandleFunc("/tasks/{taskID}", th.patchTask).Methods("PATCH")
	r.HandleFunc("/tasks/{taskID}", th.deleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{taskID}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	}).Methods("OPTIONS")
	r.HandleFunc("/tasks", th.listTasks).Methods("GET")
}
