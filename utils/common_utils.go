package utils

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func GetEnvValue(key, defaultValue string) string {
	// Load environment variables from the .env file
	godotenv.Load()

	value, ok := os.LookupEnv(key)
	if !ok || IsEmptyString(value) {
		return defaultValue
	}
	return value
}
