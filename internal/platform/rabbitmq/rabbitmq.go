package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func Connect(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *Client) Close() {
	c.ch.Close()
	c.conn.Close()
}

func (c *Client) Publish(ctx context.Context, exchange string, routingKey string, contentType string, body []byte) error {
	return c.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	})
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(queue, "", false, false, false, false, nil)
}
