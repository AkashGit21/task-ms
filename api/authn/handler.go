package authn

import (
	"github.com/AkashGit21/task-ms/internal/authn"
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	authnRouter := router.PathPrefix("/api").Subrouter()
	authn.AuthnHandler(authnRouter)

	return router, nil
}
