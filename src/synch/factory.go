package synch

import (
	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
	synch_redis "github.com/Astervia/wacraft-core/src/synch/redis"
	synch_service "github.com/Astervia/wacraft-core/src/synch/service"
)

type Factory struct {
	backend     Backend
	redisClient *synch_redis.Client
}

func NewFactory(backend Backend, redisClient *synch_redis.Client) *Factory {
	return &Factory{
		backend:     backend,
		redisClient: redisClient,
	}
}

func (f *Factory) Backend() Backend {
	return f.backend
}

func (f *Factory) RedisClient() *synch_redis.Client {
	return f.redisClient
}

// NewLock creates a DistributedLock. Due to Go generics limitations on methods,
// this is a standalone function that takes the factory as a parameter.
func NewLock[T comparable](f *Factory) synch_contract.DistributedLock[T] {
	if f.backend == BackendRedis {
		return synch_redis.NewRedisLock[T](f.redisClient)
	}
	return synch_service.NewMemoryLock[T]()
}

func (f *Factory) NewPubSub() synch_contract.PubSub {
	if f.backend == BackendRedis {
		return synch_redis.NewRedisPubSub(f.redisClient)
	}
	return synch_service.NewMemoryPubSub()
}

func (f *Factory) NewCounter() synch_contract.DistributedCounter {
	if f.backend == BackendRedis {
		return synch_redis.NewRedisCounter(f.redisClient)
	}
	return synch_service.NewMemoryCounter()
}

func (f *Factory) NewCache() synch_contract.DistributedCache {
	if f.backend == BackendRedis {
		return synch_redis.NewRedisCache(f.redisClient)
	}
	return synch_service.NewMemoryCache()
}
