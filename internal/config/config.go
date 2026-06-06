// ./internal/config/config.go

package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	YouTubeAPIKey string
}

// Load reads environment variables once and builds the config
func Load() Config {
	// load the env file into the system, so it can be used by os
	_ = godotenv.Load()


	return Config{
		Env:  getEnv("APP_ENV", "development"),
		Port: getEnv("PORT", "8080"),

		YouTubeAPIKey: getEnv("YOUTUBE_API_KEY", ""),
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

