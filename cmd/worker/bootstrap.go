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

	ticketConsumeCh, err := rabbitmq.NewChannel(rmqConn, mailer.QueueNewTicketEmail)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq ticket consume channel: %w", err)
	}
	if err := ticketConsumeCh.Qos(1, 0, false); err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq ticket qos: %w", err)
	}

	welcomeConsumeCh, err := rabbitmq.NewChannel(rmqConn, mailer.QueueWelcomeUserEmail)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq welcome consume channel: %w", err)
	}
	if err := welcomeConsumeCh.Qos(1, 0, false); err != nil {
		cleanup()
		return nil, fmt.Errorf("rabbitmq welcome qos: %w", err)
	}

	mailerWorker := mailer.NewWorker(mailerSvc, userRepo)
	ticketConsumer := mailer.NewConsumer(ticketConsumeCh, mailer.QueueNewTicketEmail, mailerWorker)
	welcomeConsumer := mailer.NewConsumer(welcomeConsumeCh, mailer.QueueWelcomeUserEmail, mailerWorker)

	slog.Info("starting worker")
	go func() {
		if err := ticketConsumer.Start(context.Background()); err != nil {
			slog.Error("ticket consumer exited with error", "error", err)
		}
	}()
	go func() {
		if err := welcomeConsumer.Start(context.Background()); err != nil {
			slog.Error("welcome consumer exited with error", "error", err)
		}
	}()
	closers = append(closers, ticketConsumer.Shutdown)
	closers = append(closers, welcomeConsumer.Shutdown)

	return cleanup, nil
}
