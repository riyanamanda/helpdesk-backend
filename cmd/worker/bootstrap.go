package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/hibiken/asynq"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
)

func bootstrap(cfg *config.Config) (func(), error) {
	var closers []func()
	cleanup := func() {
		for i := len(closers) - 1; i >= 0; i-- {
			closers[i]()
		}
	}

	slog.Info("connecting to database")
	db := database.NewPostgres(cfg.Database.ConnString())
	closers = append(closers, func() { db.Close() })

	userRepo := user.NewUserRepository(db)
	mailerSvc := mailer.NewMailerService(cfg.Email)

	redisOpt := asynq.RedisClientOpt{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 5,
		Queues:      map[string]int{"default": 1},
	})

	mux := asynq.NewServeMux()

	mailerWorker := mailer.NewWorker(mailerSvc, userRepo)
	mux.HandleFunc(mailer.TaskNewTicketEmail, mailerWorker.HandleNewTicketEmail)

	slog.Info("starting worker")
	if err := srv.Start(mux); err != nil {
		cleanup()
		return nil, fmt.Errorf("worker: %w", err)
	}
	closers = append(closers, srv.Shutdown)

	return cleanup, nil
}
