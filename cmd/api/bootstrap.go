package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	goredis "github.com/redis/go-redis/v9"

	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/cache"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/redis"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rustfs"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/storage"
	"github.com/riyanamanda/helpdesk-backend/internal/rbac"
	"github.com/riyanamanda/helpdesk-backend/internal/shared/ctxkey"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

type deps struct {
	db                *sqlx.DB
	txManager         *database.Manager
	simgosDB          *sqlx.DB
	storageService    storage.Storage
	redisClient       *goredis.Client
	cacheStore        cache.Cache
	userRepo          user.UserRepository
	outboxRepo        outbox.Repository
	permissionService ctxkey.PermissionService
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
	txManager := database.NewManager(db)

	slog.Info("running migrations")
	if err := database.RunMigrations(db); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("migrations: %w", err)
	}

	var simgosDB *sqlx.DB
	if cfg.IhsDatabase.Host != "" {
		slog.Info("connecting to simgos database")
		var simgosErr error
		simgosDB, simgosErr = database.NewMySql(cfg.IhsDatabase.MySqlConnString())
		if simgosErr != nil {
			slog.Warn("simgos database unavailable, simgos routes disabled", "error", simgosErr)
		} else {
			closers = append(closers, func() { simgosDB.Close() })
		}
	} else {
		slog.Warn("simgos database not configured, simgos routes disabled")
	}

	slog.Info("connecting to object storage")
	rustfsClient := rustfs.NewRustFSClient(cfg.Storage.Endpoint, cfg.Storage.AccessKey, cfg.Storage.SecretKey)

	if err := rustfs.InitBucket(ctx, rustfsClient, cfg.Storage.Bucket); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("rustfs bucket: %w", err)
	}

	storageService := storage.NewRustFSStorage(
		rustfsClient,
		cfg.Storage.Bucket,
	)

	slog.Info("connecting to redis")
	redisClient, err := redis.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("redis: %w", err)
	}
	closers = append(closers, func() { redisClient.Close() })

	cacheStore := cache.NewRedisCache(redisClient)
	userRepo := user.NewUserRepository(db)
	outboxRepo := outbox.NewRepository(db)
	rbacRepo := rbac.NewRBACRepository(db)
	permissionService := rbac.NewPermissionService(rbacRepo, cacheStore)

	d := &deps{
		db:                db,
		simgosDB:          simgosDB,
		storageService:    storageService,
		redisClient:       redisClient,
		cacheStore:        cacheStore,
		userRepo:          userRepo,
		outboxRepo:        outboxRepo,
		permissionService: permissionService,
		txManager:         txManager,
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.App.Host, cfg.App.Port),
		Handler: registerRoutes(cfg, d),
	}

	return server, cleanup, nil
}
