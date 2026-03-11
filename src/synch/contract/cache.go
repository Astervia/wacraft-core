package synch_contract

import "time"

type DistributedCache interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Invalidate(pattern string) error
}
