package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/riyanamanda/helpdesk-backend/internal/platform/config"
)

func NewConnection(cfg config.RabbitMQ) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: dial: %w", err)
	}
	return conn, nil
}

func NewChannel(conn *amqp.Connection, queueName string) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: open channel: %w", err)
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("rabbitmq: declare queue %q: %w", queueName, err)
	}

	return ch, nil
}
