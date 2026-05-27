package config

import (
	"os"
	"strconv"
	"time"

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
			JWTSecret:         getEnv("JWT_SECRET", "this-is-the-secret"),
			JWTExp:            getDurationEnv("JWT_EXP", 24*time.Hour),
			FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		},
		Storage: Storage{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "helpdesk-dev"),
			UseSSL:    getBoolEnv("MINIO_USE_SSL", false),
			PublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
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

func getDurationEnv(
	key string,
	fallback time.Duration,
) time.Duration {

	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}

	return d
}