package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("redis %s: %w", op, err)
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	return value, wrap("get", err)
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return wrap("set", r.client.Set(ctx, key, value, ttl).Err())
}

func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	return count, wrap("incr", err)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return wrap("del", r.client.Del(ctx, key).Err())
}

func (r *RedisCache) DeleteMany(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return wrap("del", r.client.Del(ctx, keys...).Err())
}

func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	return ttl, wrap("ttl", err)
}

func (r *RedisCache) SAdd(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	args := make([]any, len(members))
	for i, m := range members {
		args[i] = m
	}

	return wrap("sadd", r.client.SAdd(ctx, key, args...).Err())
}

func (r *RedisCache) SRem(ctx context.Context, key string, members ...string) error {
	if len(members) == 0 {
		return nil
	}

	args := make([]any, len(members))
	for i, m := range members {
		args[i] = m
	}

	return wrap("srem", r.client.SRem(ctx, key, args...).Err())
}

func (r *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	return members, wrap("smembers", err)
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return wrap("expire", r.client.Expire(ctx, key, ttl).Err())
}
