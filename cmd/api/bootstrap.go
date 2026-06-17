package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	goredis "github.com/redis/go-redis/v9"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/notification"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/firebase"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/minio"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
	"github.com/riyanamanda/helpdesk-backend/internal/user_device"
)

type deps struct {
	db                   *sqlx.DB
	ihsDB                *sqlx.DB
	storageService       storage.Storage
	redisClient          *goredis.Client
	cacheStore           cache.Cache
	userRepo             user.UserRepository
	notifier             mailer.Notifier
	notificationNotifier notification.Notifier
	permissionService    ctxkey.PermissionService
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

	slog.Info("connecting to ihs database")
	ihsDB := database.NewMySql(cfg.IhsDatabase.MySqlConnString())
	closers = append(closers, func() { ihsDB.Close() })

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

	slog.Info("connecting to rabbitmq")
	rmqConn, err := rabbitmq.NewConnection(cfg.RabbitMQ)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("rabbitmq: %w", err)
	}
	closers = append(closers, func() { rmqConn.Close() })

	publishCh, err := rabbitmq.NewChannel(rmqConn, mailer.QueueNewTicketEmail)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("rabbitmq publish channel: %w", err)
	}
	closers = append(closers, func() { publishCh.Close() })

	notifier := mailer.NewNotifier(publishCh)

	slog.Info("initializing fcm sender")
	fcmSender, err := firebase.NewFCMSender(ctx, cfg.Auth.FirebaseProjectID, cfg.Auth.FirebaseCredentialsJSON)
	if err != nil {
		slog.Warn("fcm sender unavailable, push notifications disabled", "error", err)
		fcmSender = firebase.NewNoopFCMSender()
	}

	notificationNotifier := notification.NewNotifier(
		notification.NewNotificationRepository(db),
		userRepo,
		user_device.NewUserDeviceRepository(db),
		fcmSender,
	)

	rbacRepo := rbac.NewRBACRepository(db)
	permissionService := rbac.NewPermissionService(rbacRepo, cacheStore)

	d := &deps{
		db:                   db,
		ihsDB:                ihsDB,
		storageService:       storageService,
		redisClient:          redisClient,
		cacheStore:           cacheStore,
		userRepo:             userRepo,
		notifier:             notifier,
		notificationNotifier: notificationNotifier,
		permissionService:    permissionService,
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: registerRoutes(cfg, d),
	}

	return server, cleanup, nil
}
