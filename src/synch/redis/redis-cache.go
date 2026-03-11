package synch_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *Client
}

func NewRedisCache(client *Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(key string) ([]byte, bool, error) {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("cache:" + key)

	val, err := c.client.rdb.Get(ctx, prefixed).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("redis cache get: %w", err)
	}

	return val, true, nil
}

func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("cache:" + key)

	err := c.client.rdb.Set(ctx, prefixed, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis cache set: %w", err)
	}

	return nil
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("cache:" + key)

	err := c.client.rdb.Del(ctx, prefixed).Err()
	if err != nil {
		return fmt.Errorf("redis cache delete: %w", err)
	}

	return nil
}

func (c *RedisCache) Invalidate(pattern string) error {
	ctx := context.Background()
	prefixed := c.client.PrefixKey("cache:" + pattern)

	var cursor uint64
	for {
		keys, nextCursor, err := c.client.rdb.Scan(ctx, cursor, prefixed, 100).Result()
		if err != nil {
			return fmt.Errorf("redis cache invalidate scan: %w", err)
		}

		if len(keys) > 0 {
			if err := c.client.rdb.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("redis cache invalidate delete: %w", err)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
