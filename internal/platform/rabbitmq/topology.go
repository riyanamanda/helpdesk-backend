package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Client) SetupEmailTopology() error {
	err := c.ch.ExchangeDeclare("helpdesk.events", "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare("helpdesk.email", true, false, false, false, amqp.Table{
		"x-queue-type": "quorum",
	})
	if err != nil {
		return err
	}

	err = c.ch.QueueBind("helpdesk.email", "user.created", "helpdesk.events", false, nil)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind("helpdesk.email", "ticket.created", "helpdesk.events", false, nil)
	if err != nil {
		return err
	}

	return nil
}
