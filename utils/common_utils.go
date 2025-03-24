package utils

import (
	"context"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	UserID   string `json:"userID"`
	jwt.RegisteredClaims
}

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func GetEnvValue(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || IsEmptyString(value) {
		return defaultValue
	}
	return value
}

func GetUserClaims(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value("userClaims").(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}
