package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppHost string
	AppPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("no .env loaded", "error", err)
	}

	return &Config{
		AppName: getEnv("APP_NAME", "Helpdek App"),
		AppHost: getEnv("APP_HOST", "localhost"),
		AppPort: getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
