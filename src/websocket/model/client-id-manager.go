package websocket_model

import (
	"sync"

	common_service "github.com/Astervia/wacraft-core/src/common/service"
	"github.com/google/uuid"
)

type ClientIDManager struct {
	GreatestConnID   int
	RemainingConnIDs []int // Conn ids not used between 0 and Greatest

	mu sync.Mutex
}

func (manager *ClientIDManager) CreateID(userID uuid.UUID) *ClientID {
	clientID := ClientID{
		UserID: userID,
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Add greatest connection id and return if there are no remaining connection ids
	if len(manager.RemainingConnIDs) == 0 {
		manager.GreatestConnID++
		clientID.ConnID = manager.GreatestConnID
		return &clientID
	}

	// With remaining connection ids, use the first one and remove it from the list
	clientID.ConnID = manager.RemainingConnIDs[0]
	manager.RemainingConnIDs = manager.RemainingConnIDs[1:]

	return &clientID
}

func (manager *ClientIDManager) DeleteID(clientID ClientID, deleteFromPool chan<- bool, wg *sync.WaitGroup) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	defer wg.Done()

	removedConnID := clientID.ConnID

	// If it is the greatest connection, decrement it
	if removedConnID == manager.GreatestConnID {
		manager.GreatestConnID--

		// Check if decremented value is in array. If it is, remove than decrement again and check again
		checkIfGreatestConnIDEqualsGreatestRemaining(manager)

		// If manager is empty, remove it from the pool
		if manager.GreatestConnID == -1 {
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
	common_service.InsertSorted(manager.RemainingConnIDs, removedConnID)

	manager.RemainingConnIDs = append(manager.RemainingConnIDs, clientID.ConnID)
}

func checkIfGreatestConnIDEqualsGreatestRemaining(manager *ClientIDManager) {
	remainingLen := len(manager.RemainingConnIDs)
	// If greatest remaining equals greatest conn id, remove from the array and decrement greatest conn id than check again
	if remainingLen > 0 && manager.RemainingConnIDs[remainingLen-1] == manager.GreatestConnID {
		manager.GreatestConnID--
		manager.RemainingConnIDs = manager.RemainingConnIDs[:(remainingLen - 1)]

		checkIfGreatestConnIDEqualsGreatestRemaining(manager)
	}
}
