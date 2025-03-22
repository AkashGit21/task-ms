package task

import (
	"github.com/gorilla/mux"
)

func TaskHandler(r *mux.Router) {
	r.HandleFunc("/tasks", nil).Methods("POST")
}
