package main

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

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

	dialer := func() (*amqp.Connection, error) {
		return rabbitmq.NewConnection(cfg.RabbitMQ)
	}

	mailerWorker := mailer.NewWorker(mailerSvc, userRepo)
	ticketConsumer := mailer.NewConsumer(dialer, mailer.QueueNewTicketEmail, mailerWorker)
	welcomeConsumer := mailer.NewConsumer(dialer, mailer.QueueWelcomeUserEmail, mailerWorker)

	slog.Info("starting worker")
	go func() {
		if err := ticketConsumer.Run(context.Background()); err != nil {
			slog.Error("ticket consumer exited with error", "error", err)
		}
	}()
	go func() {
		if err := welcomeConsumer.Run(context.Background()); err != nil {
			slog.Error("welcome consumer exited with error", "error", err)
		}
	}()

	return cleanup, nil
}
