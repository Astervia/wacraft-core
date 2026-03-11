package synch

type Backend string

const (
	BackendMemory Backend = "memory"
	BackendRedis  Backend = "redis"
)
