package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AkashGit21/task-ms/utils"
)

const (
	START_LOG = "START API %s %s"
	END_LOG   = "END API %s %s time=%d ms"
)

func TransactionalLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.InfoLog(fmt.Sprintf(START_LOG, r.Method, r.URL.Path))
		startTime := time.Now()

		next.ServeHTTP(w, r)

		utils.InfoLog(fmt.Sprintf(END_LOG, r.Method, r.URL.Path, time.Now().Sub(startTime).Milliseconds()))
	})
}
