package synch_service

// MemoryLock wraps MutexSwapper to implement DistributedLock[T].
type MemoryLock[T comparable] struct {
	swapper *MutexSwapper[T]
}

func NewMemoryLock[T comparable]() *MemoryLock[T] {
	return &MemoryLock[T]{
		swapper: CreateMutexSwapper[T](),
	}
}

func (m *MemoryLock[T]) Lock(key T) error {
	m.swapper.Lock(key)
	return nil
}

func (m *MemoryLock[T]) Unlock(key T) error {
	m.swapper.Unlock(key)
	return nil
}
