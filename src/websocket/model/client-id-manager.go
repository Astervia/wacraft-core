package websocket_model

import (
	"sync"

	common_service "github.com/Astervia/wacraft-core/src/common/service"
	"github.com/google/uuid"
)

type ClientIdManager struct {
	GreatestConnId   int
	RemainingConnIds []int // Conn ids not used between 0 and Greatest

	mu sync.Mutex
}

func (manager *ClientIdManager) CreateId(userId uuid.UUID) *ClientId {
	clientId := ClientId{
		UserId: userId,
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Add greatest connection id and return if there are no remaining connection ids
	if len(manager.RemainingConnIds) == 0 {
		manager.GreatestConnId++
		clientId.ConnId = manager.GreatestConnId
		return &clientId
	}

	// With remaining connection ids, use the first one and remove it from the list
	clientId.ConnId = manager.RemainingConnIds[0]
	manager.RemainingConnIds = manager.RemainingConnIds[1:]

	return &clientId
}

func (manager *ClientIdManager) DeleteId(clientId ClientId, deleteFromPool chan<- bool, wg *sync.WaitGroup) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	defer wg.Done()

	removedConnId := clientId.ConnId

	// If it is the greatest connection, decrement it
	if removedConnId == manager.GreatestConnId {
		manager.GreatestConnId--

		// Check if decremented value is in array. If it is, remove than decrement again and check again
		checkIfGreatestConnIdEqualsGreatestRemaining(manager)

		// If manager is empty, remove it from the pool
		if manager.GreatestConnId == -1 {
			// Deletion will occur only if this is the deletion of the last element
			deleteFromPool <- true
			return
		}
		deleteFromPool <- false
		return
	}

	go func() {
		deleteFromPool <- false
	}()

	// Add the connection id to the remaining inserting sorted in the array
	common_service.InsertSorted(manager.RemainingConnIds, removedConnId)

	manager.RemainingConnIds = append(manager.RemainingConnIds, clientId.ConnId)
}

func checkIfGreatestConnIdEqualsGreatestRemaining(manager *ClientIdManager) {
	remainingLen := len(manager.RemainingConnIds)
	// If greatest remaining equals greatest conn id, remove from the array and decrement greatest conn id than check again
	if remainingLen > 0 && manager.RemainingConnIds[remainingLen-1] == manager.GreatestConnId {
		manager.GreatestConnId--
		manager.RemainingConnIds = manager.RemainingConnIds[:(remainingLen - 1)]

		checkIfGreatestConnIdEqualsGreatestRemaining(manager)
	}
}
