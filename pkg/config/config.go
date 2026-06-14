package config

import (
	"os"
)

type Config struct {
	Port       string
	BaseURL    string
	APISecret  string
	StorageDir string
	LogLevel   string
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		BaseURL:    getEnv("BASE_URL", "https://api.zynu.net"),
		APISecret:  getEnv("API_SECRET", ""),
		StorageDir: getEnv("STORAGE_DIR", "./storage"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
