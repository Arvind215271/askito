// ./internal/config/config
package config

import (
	"os"
)

type Config struct {
	Env  string
	Port string
}

// Load reads environment variables once and builds the config
func Load() Config {
	return Config{
		Env:  getEnv("APP_ENV", "development"),
		Port: getEnv("PORT", "8080"),
	}
}

// getEnv returns fallback if env is not set
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}