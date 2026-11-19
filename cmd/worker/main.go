package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/user"
	"github.com/riyanamanda/helpdesk-backend/internal/worker"
)

func main() {
	if err := run(); err != nil {
		slog.Error("worker stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	if names := cfg.InsecureDefaults(); len(names) > 0 {
		slog.Warn("insecure default config in use", "env", names)
	}
	if cfg.Email.Host == "" {
		slog.Warn("mail not configured, emails will fail")
	}

	// RabbitMQ Initial
	client, err := rabbitmq.Connect(cfg.RabbitMQ.RabbitMQConnString())
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	defer client.Close()

	// Topology -- exchange, queue, bind
	err = client.SetupEmailTopology()
	if err != nil {
		return fmt.Errorf("setup rabbitmq topology: %w", err)
	}

	// Database connection
	db, err := database.NewPostgres(cfg.Database.ConnString())
	if err != nil {
		return err
	}
	defer db.Close()

	// Dependencies
	outboxRepo := outbox.NewRepository(db)
	userRepo := user.NewRepository(db)
	mailer := mailer.NewMailer(cfg.Email)
	publisher := worker.NewPublisher(outboxRepo, client)

	// goroutine
	go publisher.Run(ctx)

	queues := []string{
		rabbitmq.QueueWelcome,
		rabbitmq.QueueTicket,
		rabbitmq.QueueTicketActivity,
		rabbitmq.QueuePasswordReset,
	}

	errCh := make(chan error, len(queues))
	for _, queue := range queues {
		consumer := worker.NewConsumer(client, mailer, *userRepo, queue, cfg.App.URL)
		go func() {
			if err := consumer.Run(ctx); err != nil {
				errCh <- fmt.Errorf("%s consumer: %w", queue, err)
			}
		}()
	}

	slog.Info("worker started")

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutdown signal received")

	return nil
}
