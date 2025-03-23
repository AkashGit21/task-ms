package middleware

import (
	"net/http"

	"github.com/AkashGit21/task-ms/utils"
)

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Handle unknown panic
				utils.InfoLog("Panic recovered: ", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
