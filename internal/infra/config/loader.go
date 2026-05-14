package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		panic(".env files not loaded properly")
	}

	return &Config{
		App: App{
			Name: getEnv("APP_NAME", "Helpdesk Api"),
			Host: getEnv("APP_HOST", "localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "helpdesk_db"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: Auth{
			JWTSecret:            getEnv("JWT_SECRET", "this-is-the-secret"),
			JWTExpirationMinutes: getIntEnv("JWT_EXP", getIntEnv("JWT_EXPIRED", 60)),
		},
		Storage: Storage{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINOI_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINOI_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINOI_BUCKET", "helpdesk-dev"),
			UseSSL:    getBoolEnv("MINIO_USE_SSL", false),
			PublicURL: getEnv("MINOI_PUBLIC_URL", "http://localhost:9000"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed
		}
	}

	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return parsed
		}
	}

	return fallback
}
