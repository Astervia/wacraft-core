package synch_contract

type DistributedLock[T comparable] interface {
	Lock(key T) error
	Unlock(key T) error
	// TryLock attempts to acquire the lock without blocking.
	// Returns (true, nil) if acquired, (false, nil) if already held, or (false, err) on error.
	TryLock(key T) (bool, error)
}
