package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppHost string
	AppPort string

	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPass    string
	DBSSLMode string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("no .env loaded", "error", err)
	}

	return &Config{
		AppName: getEnv("APP_NAME", "Helpdek App"),
		AppHost: getEnv("APP_HOST", "localhost"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBName:    getEnv("DB_NAME", "helpdesk_db"),
		DBUser:    getEnv("DB_USERNAME", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPass,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
