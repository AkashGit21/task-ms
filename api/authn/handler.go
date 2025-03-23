package authn

import (
	"github.com/AkashGit21/task-ms/internal/authn"
	middlewares "github.com/AkashGit21/task-ms/internal/middleware"
	"github.com/gorilla/mux"
)

func New() (*mux.Router, error) {
	router := mux.NewRouter()

	authnRouter := router.PathPrefix("/api").Subrouter()
	authnRouter.Use(middlewares.TransactionalLogMiddleware)
	authnRouter.Use(middlewares.PanicRecoveryMiddleware)
	authn.AuthnHandler(authnRouter)

	return router, nil
}
