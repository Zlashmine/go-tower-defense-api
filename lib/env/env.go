package env

import (
	"os"
	"strconv"
)

func GetString(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return defaultValue
	}

	return val
}

func GetInt(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		return defaultValue
	}

	valInt, err := strconv.Atoi(val)

	if err != nil {
		return defaultValue
	}

	return valInt
}

func GetBool(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)

	if !ok {
		return defaultValue
	}

	valBool, err := strconv.ParseBool(val)

	if err != nil {
		return defaultValue
	}

	return valBool
}
