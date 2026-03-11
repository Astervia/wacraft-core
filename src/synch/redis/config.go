package synch_redis

import "time"

type Config struct {
	URL       string
	Password  string
	DB        int
	KeyPrefix string
	LockTTL   time.Duration
	CacheTTL  time.Duration
}
