package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultJWTSecret        = "this-is-the-secret"
	defaultDBPassword       = "postgres"
	defaultStorageKey       = "rustfsadmin"
	defaultRabbitMQPassword = "guest"
)

func (c *Config) InsecureDefaults() []string {
	var names []string

	if c.Auth.JWTSecret == defaultJWTSecret {
		names = append(names, "JWT_SECRET")
	}
	if c.Database.Password == defaultDBPassword {
		names = append(names, "DB_PASSWORD")
	}
	if c.Storage.AccessKey == defaultStorageKey {
		names = append(names, "RUSTFS_ACCESS_KEY")
	}
	if c.Storage.SecretKey == defaultStorageKey {
		names = append(names, "RUSTFS_SECRET_KEY")
	}
	if c.RabbitMQ.Password == defaultRabbitMQPassword {
		names = append(names, "RABBITMQ_PASSWORD")
	}

	return names
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, reading from environment")
	}

	return &Config{
		App: App{
			Name:        getEnv("APP_NAME", "Helpdesk Api"),
			Host:        getEnv("APP_HOST", "localhost"),
			Port:        getEnv("APP_PORT", "8080"),
			URL:         getEnv("APP_URL", "http://localhost:3000"),
			CORSOrigins: getStringSliceEnv("CORS_ORIGINS", []string{"http://localhost:3000"}),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "helpdesk_db"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", defaultDBPassword),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		IhsDatabase: Database{
			Host:     getEnv("IHS_DB_HOST", ""),
			Port:     getEnv("IHS_DB_PORT", "3306"),
			Name:     getEnv("IHS_DB_NAME", "kemkes-ihs"),
			Username: getEnv("IHS_DB_USERNAME", "root"),
			Password: getEnv("IHS_DB_PASSWORD", ""),
		},
		Auth: Auth{
			JWTSecret:         getEnv("JWT_SECRET", defaultJWTSecret),
			JWTExp:            getDurationEnv("JWT_EXP", 72*time.Hour),
			FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		},
		Storage: Storage{
			Endpoint:  getEnv("RUSTFS_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("RUSTFS_ACCESS_KEY", defaultStorageKey),
			SecretKey: getEnv("RUSTFS_SECRET_KEY", defaultStorageKey),
			Bucket:    getEnv("RUSTFS_BUCKET", "helpdesk-dev"),
		},
		Redis: Redis{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		RabbitMQ: RabbitMQ{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			Username: getEnv("RABBITMQ_USERNAME", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", defaultRabbitMQPassword),
		},
		Email: Email{
			Host:     getEnv("MAIL_HOST", ""),
			Port:     getEnv("MAIL_PORT", ""),
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", ""),
			UseSSL:   getBoolEnv("MAIL_SSL", false),
		},
		Antrol: Antrol{
			Domain:   getEnv("ANTROL_DOMAIN", ""),
			Username: getEnv("ANTROL_USERNAME", ""),
			Password: getEnv("ANTROL_PASSWORD", ""),
		},
		BPJS: BPJS{
			ConsID:     getEnv("CONS_ID", ""),
			ConsSecret: getEnv("CONS_SECRET", ""),
			VClaimURL:  getEnv("VCLAIM_URL", ""),
			VClaimKey:  getEnv("VCLAIM_KEY", ""),
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
		slog.Warn("invalid env value, using default", "env", key, "default", fallback)
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
		slog.Warn("invalid env value, using default", "env", key, "default", fallback.String())
		return fallback
	}

	return d
}
