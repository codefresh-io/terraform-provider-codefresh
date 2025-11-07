// Package envutil provides utility functions for working with environment variables.
package envutil

import (
	"os"
	"strconv"
)

// GetEnvAsString retrieves an environment variable as a string, returning the default value if not set
func GetEnvAsString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retrieves an environment variable as an integer, returning the default value if not set or invalid
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool retrieves an environment variable as a boolean, returning the default value if not set or invalid
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
