package env

import (
	"os"
	"strconv"
)

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err != nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
