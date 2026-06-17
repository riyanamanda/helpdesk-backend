package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, reading from environment")
	}

	return &Config{
		App: App{
			Name:        getEnv("APP_NAME", "Helpdesk Api"),
			Host:        getEnv("APP_HOST", "localhost"),
			Port:        getEnv("APP_PORT", "8080"),
			CORSOrigins: getStringSliceEnv("CORS_ORIGINS", []string{"http://localhost:3000"}),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "helpdesk_db"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		IhsDatabase: Database{
			Host:     getEnv("IHS_DB_HOST", "localhost"),
			Port:     getEnv("IHS_DB_PORT", "3306"),
			Name:     getEnv("IHS_DB_NAME", "kemkes-ihs"),
			Username: getEnv("IHS_DB_USERNAME", "root"),
			Password: getEnv("IHS_DB_PASSWORD", ""),
		},
		Auth: Auth{
			JWTSecret:               getEnv("JWT_SECRET", "this-is-the-secret"),
			JWTExp:                  getDurationEnv("JWT_EXP", 24*time.Hour),
			FirebaseProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
			FirebaseCredentialsJSON: getEnv("FIREBASE_CREDENTIALS_JSON", ""),
		},
		Storage: Storage{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "helpdesk-dev"),
			UseSSL:    getBoolEnv("MINIO_USE_SSL", false),
		},
		Redis: Redis{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Email: Email{
			Host:     getEnv("MAIL_HOST", ""),
			Port:     getEnv("MAIL_PORT", ""),
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", ""),
			UseSSL:   getBoolEnv("MAIL_SSL", false),
		},
		RabbitMQ: RabbitMQ{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			Username: getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
		},
	}
}

func getStringSliceEnv(key string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	parts := strings.Split(v, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}

	return result
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

func getDurationEnv(key string, fallback time.Duration) time.Duration {
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
