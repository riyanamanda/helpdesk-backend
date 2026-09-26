package main

import (
	"context"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	slog.Info("Worker started...")

	// RabbitMQ Initial
	client, err := rabbitmq.Connect(cfg.RabbitMQ.RabbitMQConnString())
	if err != nil {
		slog.Error("RabbitMQ connect failed", "error", err)
		return
	}
	defer client.Close()

	// Topology -- exchange, queue, bind
	err = client.SetupEmailTopology()
	if err != nil {
		slog.Error("RabbitMQ setup email topology failed", "error", err)
		return
	}

	// Database connection
	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	// Dependencies
	outboxRepo := outbox.NewRepository(db)
	userRepo := user.NewUserRepository(db)
	mailer := mailer.NewMailer(cfg.Email)
	publisher := worker.NewPublisher(outboxRepo, client)
	consumer := worker.NewConsumer(client, mailer, userRepo)

	// goroutine
	go publisher.Run(ctx)

	go func() {
		if err := consumer.Run(ctx); err != nil {
			slog.Error("consumer stopped", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")
}
