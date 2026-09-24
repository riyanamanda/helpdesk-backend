package main

import (
	"context"
	"log/slog"

	"github.com/riyanamanda/helpdesk-backend/internal/mailer"
	"github.com/riyanamanda/helpdesk-backend/internal/outbox"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/database"
	"github.com/riyanamanda/helpdesk-backend/internal/platform/rabbitmq"
	"github.com/riyanamanda/helpdesk-backend/internal/worker"
)

type UserCreatedEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	ctx := context.Background()
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
	}

	db := database.NewPostgres(cfg.Database.ConnString())
	defer db.Close()

	outboxRepo := outbox.NewRepository(db)
	mailer := mailer.NewMailer(cfg.Email)
	publisher := worker.NewPublisher(outboxRepo, client)
	consumer := worker.NewConsumer(client, mailer)

	// goroutine
	go publisher.Run(ctx)
	go func() {
		if err := consumer.Run(ctx); err != nil {
			slog.Error("consumer stopped", "error", err)
		}
	}()

	select {}
}
