package websocket_model

import (
	"sync"

	"github.com/google/uuid"
)

type ClientPool struct {
	mu       *sync.Mutex
	managers map[uuid.UUID]*ClientIDManager
}

func (p *ClientPool) CreateID(userID uuid.UUID) *ClientID {
	p.mu.Lock()
	defer p.mu.Unlock()

	manager, ok := (p.managers)[userID]

	// Add manager if it doesn't exit
	if !ok {
		manager = &ClientIDManager{
			GreatestConnID:   -1,
			RemainingConnIDs: []int{},
		}
		(p.managers)[userID] = manager
	}

	return manager.CreateID(userID)
}

func (p *ClientPool) DeleteID(clientID ClientID) {
	deleteFromPool := make(chan bool)
	var wg sync.WaitGroup
	defer wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()
	manager, ok := (p.managers)[clientID.UserID]
	if !ok {
		return
	}

	wg.Add(1)
	go manager.DeleteID(clientID, deleteFromPool, &wg)

	if <-deleteFromPool {
		delete(p.managers, clientID.UserID)
	}
}

func (p *ClientPool) DeleteManager(userID uuid.UUID) {
	p.mu.Lock()
	delete(p.managers, userID)
	p.mu.Unlock()
}

func CreateClientPool() *ClientPool {
	var mu sync.Mutex
	return &ClientPool{
		mu:       &mu,
		managers: make(map[uuid.UUID]*ClientIDManager),
	}
}
