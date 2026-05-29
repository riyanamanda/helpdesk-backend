package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"

	"github.com/riyanamanda/helpdesk-backend/internal/infra/config"
)

func NewRedisClient(ctx context.Context, cfg config.Redis) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{

		Addr: net.JoinHostPort(cfg.Host, cfg.Port),

		Password: cfg.Password,

		DB: 0,
	})

	if err := client.Ping(ctx).Err(); err != nil {

		return nil, fmt.Errorf("failed to connect redis: %w", err)

	}

	return client, nil

}
