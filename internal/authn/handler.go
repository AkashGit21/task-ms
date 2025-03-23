package authn

import (
	"github.com/AkashGit21/task-ms/lib/persistence/mysql"
	"github.com/gorilla/mux"
)

type authnHandler struct {
	mysql.AuthnOps
}

func NewAuthnAPIHandler() *authnHandler {
	persistenceDB, err := mysql.NewUserPersistenceLayer()
	if err != nil {
		panic(err)
	}

	return &authnHandler{
		persistenceDB,
	}
}

func AuthnHandler(r *mux.Router) {
	anh := NewAuthnAPIHandler()

	r.HandleFunc("/tasks", anh.IsAuthenticated).Methods("POST")
}
