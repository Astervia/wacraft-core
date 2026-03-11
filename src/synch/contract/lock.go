package synch_contract

type DistributedLock[T comparable] interface {
	Lock(key T) error
	Unlock(key T) error
}
