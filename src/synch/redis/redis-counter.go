package synch_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCounter struct {
	client *Client
}

func NewRedisCounter(client *Client) *RedisCounter {
	return &RedisCounter{client: client}
}

func (c *RedisCounter) Increment(key string, delta int64) (int64, error) {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("counter:" + key)

	val, err := c.client.rdb.IncrBy(ctx, prefixed, delta).Result()
	if err != nil {
		return 0, fmt.Errorf("redis increment: %w", err)
	}

	return val, nil
}

func (c *RedisCounter) Get(key string) (int64, error) {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("counter:" + key)

	val, err := c.client.rdb.Get(ctx, prefixed).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("redis get counter: %w", err)
	}

	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("redis parse counter: %w", err)
	}

	return n, nil
}

func (c *RedisCounter) SetTTL(key string, ttl time.Duration) error {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("counter:" + key)

	err := c.client.rdb.Expire(ctx, prefixed, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set counter TTL: %w", err)
	}

	return nil
}

func (c *RedisCounter) Delete(key string) error {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("counter:" + key)

	err := c.client.rdb.Del(ctx, prefixed).Err()
	if err != nil {
		return fmt.Errorf("redis delete counter: %w", err)
	}

	return nil
}
