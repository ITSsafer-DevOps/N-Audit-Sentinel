package config

import "os"

// GetEnv returns the value for key or fallback if not set
func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
