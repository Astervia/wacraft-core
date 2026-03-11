package synch_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb    *redis.Client
	config Config
}

func NewClient(config Config) (*Client, error) {
	opts, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if config.Password != "" {
		opts.Password = config.Password
	}
	opts.DB = config.DB

	rdb := redis.NewClient(opts)

	client := &Client{
		rdb:    rdb,
		config: config,
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) PingWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.Ping(ctx)
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) PrefixKey(key string) string {
	return c.config.KeyPrefix + key
}

func (c *Client) Redis() *redis.Client {
	return c.rdb
}

func (c *Client) Config() Config {
	return c.config
}
