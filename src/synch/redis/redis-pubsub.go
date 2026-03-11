package synch_redis

import (
	"context"
	"fmt"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	client *Client
}

type redisSubscription struct {
	pubsub *redis.PubSub
	ch     chan []byte
	cancel context.CancelFunc
}

func NewRedisPubSub(client *Client) *RedisPubSub {
	return &RedisPubSub{client: client}
}

func (p *RedisPubSub) Publish(channel string, message []byte) error {
	ctx := context.Background()
	prefixed := p.client.PrefixKey(channel)
	err := p.client.rdb.Publish(ctx, prefixed, message).Err()
	if err != nil {
		return fmt.Errorf("redis publish: %w", err)
	}
	return nil
}

func (p *RedisPubSub) Subscribe(channel string) (synch_contract.Subscription, error) {
	ctx, cancel := context.WithCancel(context.Background())
	prefixed := p.client.PrefixKey(channel)
	pubsub := p.client.rdb.Subscribe(ctx, prefixed)

	// Wait for confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		cancel()
		pubsub.Close()
		return nil, fmt.Errorf("redis subscribe: %w", err)
	}

	ch := make(chan []byte, 100)
	sub := &redisSubscription{
		pubsub: pubsub,
		ch:     ch,
		cancel: cancel,
	}

	// Forward Redis messages to the channel
	go func() {
		redisCh := pubsub.Channel()
		for msg := range redisCh {
			select {
			case ch <- []byte(msg.Payload):
			default:
				// Drop if buffer full
			}
		}
		close(ch)
	}()

	return sub, nil
}

func (s *redisSubscription) Channel() <-chan []byte {
	return s.ch
}

func (s *redisSubscription) Unsubscribe() error {
	s.cancel()
	return s.pubsub.Close()
}
