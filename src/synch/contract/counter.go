package synch_contract

import "time"

type DistributedCounter interface {
	Increment(key string, delta int64) (int64, error)
	Get(key string) (int64, error)
	SetTTL(key string, ttl time.Duration) error
	Delete(key string) error
}
