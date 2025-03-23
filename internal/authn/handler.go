package authn

import (
	"github.com/AkashGit21/task-ms/lib/persistence/mysql"
	"github.com/gorilla/mux"
)

type authnHandler struct {
	mysql.AuthnOps
}

func newAuthnAPIHandler() *authnHandler {
	persistenceDB, err := mysql.NewUserPersistenceLayer()
	if err != nil {
		panic(err)
	}

	return &authnHandler{
		persistenceDB,
	}
}

func AuthnHandler(r *mux.Router) {
	anh := newAuthnAPIHandler()

	r.HandleFunc("/login", anh.IsAuthenticated).Methods("POST")
}
