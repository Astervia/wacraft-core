package websocket_model

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientId struct {
	UserId uuid.UUID
	ConnId int
}

func (c *ClientId) String() string {
	return fmt.Sprintf("%s-%d", c.UserId.String(), c.ConnId)
}

func CompareClients(client1, client2 ClientId) bool {
	return client1.UserId == client2.UserId && client1.ConnId == client2.ConnId
}
