package websocket_model

import (
	"fmt"
	"sync"

	"github.com/pterm/pterm"
)

type Channel[T, U any, V comparable] struct {
	Clients      map[V]Client[T]
	ClientsMutex *sync.Mutex
}

func (c *Channel[T, U, V]) AppendClient(client Client[T], key V) {
	c.ClientsMutex.Lock()
	defer c.ClientsMutex.Unlock()
	c.Clients[key] = client
}

func (c *Channel[T, U, V]) RemoveClient(key V) {
	c.ClientsMutex.Lock()
	defer c.ClientsMutex.Unlock()

	_, ok := c.Clients[key]
	if !ok {
		return
	}

	delete(c.Clients, key)
}

func (c *Channel[T, U, V]) BroadcastJsonMultithread(data U) {
	var wg sync.WaitGroup

	for i := range c.Clients {
		wg.Add(1)

		go func(i V) {
			defer wg.Done()
			c.ClientsMutex.Lock()
			defer c.ClientsMutex.Unlock()

			client, ok := c.Clients[i]
			if !ok {
				return
			}

			if err := client.Connection.WriteJSON(data); err != nil {
				pterm.DefaultLogger.Info(fmt.Sprintf("Error sending message to client %v", err))
			}
		}(i)
	}

	wg.Wait()
}

func (c *Channel[T, U, V]) BroadcastMessageMultithread(
	messageType int,
	data []byte,
) {
	var wg sync.WaitGroup

	for i := range c.Clients {
		wg.Add(1)

		go func(i V) {
			defer wg.Done()
			c.ClientsMutex.Lock()
			defer c.ClientsMutex.Unlock()

			client, ok := c.Clients[i]
			if !ok {
				return
			}

			if err := client.Connection.WriteMessage(
				messageType,
				data,
			); err != nil {
				pterm.DefaultLogger.Info(fmt.Sprintf("Error sending message to client %v", err))
			}
		}(i)
	}

	wg.Wait()
}

func CreateChannel[T, U any, V comparable]() *Channel[T, U, V] {
	var clientsMutex sync.Mutex
	return &Channel[T, U, V]{
		ClientsMutex: &clientsMutex,
		Clients:      make(map[V]Client[T]),
	}
}
