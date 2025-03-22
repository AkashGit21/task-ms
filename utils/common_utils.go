package utils

import (
	"os"
	"strings"
)

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
