package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
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

	slog.Info("connecting to rabbitmq")
	rmqConn, err := rabbitmq.NewConnection(cfg.RabbitMQ)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq: %w", err)
	}
	closers = append(closers, func() { rmqConn.Close() })

	consumeCh, err := rabbitmq.NewChannel(rmqConn, mailer.QueueNewTicketEmail)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq consume channel: %w", err)
	}
	if err := consumeCh.Qos(1, 0, false); err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq qos: %w", err)
	}

	mailerWorker := mailer.NewWorker(mailerSvc, userRepo)
	consumer := mailer.NewConsumer(consumeCh, mailerWorker)

	slog.Info("starting worker")
	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			slog.Error("consumer exited with error", "error", err)
		}
	}()
	closers = append(closers, consumer.Shutdown)

	return cleanup, nil
}
