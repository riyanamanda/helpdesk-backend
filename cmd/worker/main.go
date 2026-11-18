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

	slog.Info("Worker started...")

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
	welcomeConsumer := worker.NewConsumer(client, mailer, *userRepo, rabbitmq.QueueWelcome)
	ticketConsumer := worker.NewConsumer(client, mailer, *userRepo, rabbitmq.QueueTicket)
	passwordResetConsumer := worker.NewConsumer(client, mailer, *userRepo, rabbitmq.QueuePasswordReset)

	// goroutine
	go publisher.Run(ctx)

	errCh := make(chan error, 3)
	go func() {
		if err := welcomeConsumer.Run(ctx); err != nil {
			errCh <- fmt.Errorf("welcome consumer: %w", err)
		}
	}()

	go func() {
		if err := ticketConsumer.Run(ctx); err != nil {
			errCh <- fmt.Errorf("ticket consumer: %w", err)
		}
	}()

	go func() {
		if err := passwordResetConsumer.Run(ctx); err != nil {
			errCh <- fmt.Errorf("password reset consumer: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutdown signal received")

	return nil
}
