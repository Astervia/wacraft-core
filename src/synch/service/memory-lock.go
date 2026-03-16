package synch_service

import "sync"

// MemoryLock wraps MutexSwapper to implement DistributedLock[T].
type MemoryLock[T comparable] struct {
	swapper *MutexSwapper[T]
	tryHeld sync.Map // tracks keys acquired via TryLock
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

func (m *MemoryLock[T]) TryLock(key T) (bool, error) {
	_, loaded := m.tryHeld.LoadOrStore(key, struct{}{})
	return !loaded, nil
}

func (m *MemoryLock[T]) Unlock(key T) error {
	// Clean up either a TryLock or a regular Lock entry (one of these is always a no-op).
	m.tryHeld.Delete(key)
	m.swapper.Unlock(key)
	return nil
}
