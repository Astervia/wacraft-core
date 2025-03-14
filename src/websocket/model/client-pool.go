package websocket_model

import (
	"sync"

	"github.com/google/uuid"
)

type ClientPool struct {
	mu       *sync.Mutex
	managers map[uuid.UUID]*ClientIdManager
}

func (p *ClientPool) CreateId(userId uuid.UUID) *ClientId {
	p.mu.Lock()
	defer p.mu.Unlock()

	manager, ok := (p.managers)[userId]

	// Add manager if it doesn't exit
	if !ok {
		manager = &ClientIdManager{
			GreatestConnId:   -1,
			RemainingConnIds: []int{},
		}
		(p.managers)[userId] = manager
	}

	return manager.CreateId(userId)
}

func (p *ClientPool) DeleteId(clientId ClientId) {
	deleteFromPool := make(chan bool)
	var wg sync.WaitGroup
	defer wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()
	manager, ok := (p.managers)[clientId.UserId]
	if !ok {
		return
	}

	wg.Add(1)
	go manager.DeleteId(clientId, deleteFromPool, &wg)

	if <-deleteFromPool {
		delete(p.managers, clientId.UserId)
	}
}

func (p *ClientPool) DeleteManager(userId uuid.UUID) {
	p.mu.Lock()
	delete(p.managers, userId)
	p.mu.Unlock()
}

func CreateClientPool() *ClientPool {
	var mu sync.Mutex
	return &ClientPool{
		mu:       &mu,
		managers: make(map[uuid.UUID]*ClientIdManager),
	}
}
