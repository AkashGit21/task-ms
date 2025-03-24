package authn

import (
	"github.com/AkashGit21/task-ms/lib/persistence/mysql"
	"github.com/AkashGit21/task-ms/utils"
	"github.com/gorilla/mux"
)

type authnHandler struct {
	mysql.AuthnOps
	JWTSecret []byte
}

func newAuthnAPIHandler() *authnHandler {
	persistenceDB, err := mysql.NewUserPersistenceLayer()
	if err != nil {
		panic(err)
	}

	utils.ErrorLog("JWT secret ", utils.GetEnvValue("AUTH_JWT_SECRET", "DEFAULT SECRET"))

	return &authnHandler{
		persistenceDB,
		[]byte(utils.GetEnvValue("AUTH_JWT_SECRET", "DEFAULT SECRET")),
	}
}

func AuthnHandler(r *mux.Router) {
	anh := newAuthnAPIHandler()

	r.HandleFunc("/login", anh.UserLogin).Methods("POST")
}
