package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/riyanamanda/helpdesk-backend/internal/event"
)

func (c *Client) SetupEmailTopology() error {
	// exchange
	if err := c.ch.ExchangeDeclare(ExchangeEvent, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	// welcome DLX
	if err := c.ch.ExchangeDeclare(ExchangeWelcomeDLX, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// ticket DLX
	if err := c.ch.ExchangeDeclare(ExchangeTicketDLX, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// welcome queue
	if _, err := c.ch.QueueDeclare(QueueWelcome, true, false, false, false, amqp.Table{
		"x-queue-type":              "quorum",
		"x-dead-letter-exchange":    ExchangeWelcomeDLX,
		"x-dead-letter-routing-key": event.UserCreated,
	}); err != nil {
		return err
	}

	// ticket queue
	if _, err := c.ch.QueueDeclare(QueueTicket, true, false, false, false, amqp.Table{
		"x-queue-type":              "quorum",
		"x-dead-letter-exchange":    ExchangeTicketDLX,
		"x-dead-letter-routing-key": event.TicketCreated,
	}); err != nil {
		return err
	}

	// password reset queue | fire and forget
	if _, err := c.ch.QueueDeclare(QueuePasswordReset, true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	}); err != nil {
		return err
	}

	// welcome DLQ
	if _, err := c.ch.QueueDeclare(QueueWelcomeDLQ, true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	}); err != nil {
		return err
	}

	// ticket DLQ
	if _, err := c.ch.QueueDeclare(QueueTicketDLQ, true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	}); err != nil {
		return err
	}
	// end queue

	// main queue
	if err := c.ch.QueueBind(QueueWelcome, event.UserCreated, ExchangeEvent, false, nil); err != nil {
		return err
	}

	if err := c.ch.QueueBind(QueueTicket, event.TicketCreated, ExchangeEvent, false, nil); err != nil {
		return err
	}

	if err := c.ch.QueueBind(QueuePasswordReset, event.PasswordResetRequested, ExchangeEvent, false, nil); err != nil {
		return err
	}
	// end main queue

	// DLQ queue
	if err := c.ch.QueueBind(QueueWelcomeDLQ, event.UserCreated, ExchangeWelcomeDLX, false, nil); err != nil {
		return err
	}

	if err := c.ch.QueueBind(QueueTicketDLQ, event.TicketCreated, ExchangeTicketDLX, false, nil); err != nil {
		return err
	}
	// End DLQ queue

	return nil
}
