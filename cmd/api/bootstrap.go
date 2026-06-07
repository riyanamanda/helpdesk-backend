package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	goredis "github.com/redis/go-redis/v9"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/minio"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type deps struct {
	db             *sqlx.DB
	storageService storage.Storage
	redisClient    *goredis.Client
	cacheStore     cache.Cache
	userRepo       user.UserRepository
	notifier       mailer.Notifier
}

func bootstrap(ctx context.Context, cfg *config.Config) (*http.Server, func(), error) {
	var closers []func()
	cleanup := func() {
		for i := len(closers) - 1; i >= 0; i-- {
			closers[i]()
		}
	}

	slog.Info("connecting to database")
	db := database.NewPostgres(cfg.Database.ConnString())
	closers = append(closers, func() { db.Close() })

	slog.Info("running migrations")
	if err := database.RunMigrations(db); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("migrations: %w", err)
	}

	slog.Info("connecting to minio")
	minioClient, err := minio.NewMinioClient(
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.UseSSL,
	)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("minio: %w", err)
	}

	if err := minio.InitBucket(ctx, minioClient, cfg.Storage.Bucket); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("minio bucket: %w", err)
	}

	storageService := storage.NewMinioStorage(minioClient, cfg.Storage.Bucket)

	slog.Info("connecting to redis")
	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("redis: %w", err)
	}
	closers = append(closers, func() { redisClient.Close() })

	cacheStore := cache.NewRedisCache(redisClient)
	userRepo := user.NewUserRepository(db)

	asynqOpt := asynq.RedisClientOpt{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}
	asynqClient := asynq.NewClient(asynqOpt)
	closers = append(closers, func() { asynqClient.Close() })

	notifier := mailer.NewNotifier(asynqClient)

	d := &deps{
		db:             db,
		storageService: storageService,
		redisClient:    redisClient,
		cacheStore:     cacheStore,
		userRepo:       userRepo,
		notifier:       notifier,
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: registerRoutes(cfg, d),
	}

	return server, cleanup, nil
}
