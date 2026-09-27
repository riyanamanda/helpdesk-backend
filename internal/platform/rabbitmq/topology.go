package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
)

func (c *Client) SetupEmailTopology() error {
	// exchange
	err := c.ch.ExchangeDeclare(ExchangeEvent, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// queue
	_, err = c.ch.QueueDeclare(QueueWelcome, true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	})
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(QueueTicket, true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	})
	if err != nil {
		return err
	}
	// end queue

	// bind
	err = c.ch.QueueBind(QueueWelcome, event.UserCreated, ExchangeEvent, false, nil)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(QueueTicket, event.TicketCreated, ExchangeEvent, false, nil)
	if err != nil {
		return err
	}
	// end bind

	return nil
}
