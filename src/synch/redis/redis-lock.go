package synch_redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// unlockScript atomically checks and deletes the lock only if the owner matches.
var unlockScript = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
end
return 0
`)

type RedisLock[T comparable] struct {
	client   *Client
	ttl      time.Duration
	ownerID  string
	owners   map[T]string
	ownersMu sync.Mutex
}

func NewRedisLock[T comparable](client *Client) *RedisLock[T] {
	return &RedisLock[T]{
		client:  client,
		ttl:     client.config.LockTTL,
		ownerID: uuid.New().String(),
		owners:  make(map[T]string),
	}
}

func (l *RedisLock[T]) TryLock(key T) (bool, error) {
	ctx := context.Background()
	redisKey := l.client.PrefixKey(fmt.Sprintf("lock:%v", key))
	value := fmt.Sprintf("%s:%s", l.ownerID, uuid.New().String())

	ok, err := l.client.rdb.SetNX(ctx, redisKey, value, l.ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis try lock: %w", err)
	}

	if ok {
		l.ownersMu.Lock()
		l.owners[key] = value
		l.ownersMu.Unlock()
	}

	return ok, nil
}

func (l *RedisLock[T]) Lock(key T) error {
	ctx := context.Background()
	redisKey := l.client.PrefixKey(fmt.Sprintf("lock:%v", key))
	value := fmt.Sprintf("%s:%s", l.ownerID, uuid.New().String())

	for {
		ok, err := l.client.rdb.SetNX(ctx, redisKey, value, l.ttl).Result()
		if err != nil {
			return fmt.Errorf("redis lock: %w", err)
		}

		if ok {
			l.ownersMu.Lock()
			l.owners[key] = value
			l.ownersMu.Unlock()
			return nil
		}

		// Wait briefly before retrying
		time.Sleep(10 * time.Millisecond)
	}
}

func (l *RedisLock[T]) Unlock(key T) error {
	l.ownersMu.Lock()
	value, exists := l.owners[key]
	if exists {
		delete(l.owners, key)
	}
	l.ownersMu.Unlock()

	if !exists {
		return nil
	}

	ctx := context.Background()
	redisKey := l.client.PrefixKey(fmt.Sprintf("lock:%v", key))

	_, err := unlockScript.Run(ctx, l.client.rdb, []string{redisKey}, value).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("redis unlock: %w", err)
	}

	return nil
}
