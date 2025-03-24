package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/AkashGit21/task-ms/utils"
	"github.com/golang-jwt/jwt/v5"
)

/** Adds Bearer JWT token authentication with user details **/
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		jwtSecret := utils.GetEnvValue("AUTH_JWT_SECRET", "some secret")

		token, err := jwt.ParseWithClaims(tokenString, &utils.UserClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
		utils.InfoLog("middleware jwt secret", []byte(jwtSecret))

		if err != nil {
			utils.ErrorLog("Token invalid", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*utils.UserClaims)
		if ok && token.Valid {
			// check if the token expires or not
			if float64(time.Now().UTC().Unix()) > float64(claims.ExpiresAt.UTC().Unix()) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user claims to the request context
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
