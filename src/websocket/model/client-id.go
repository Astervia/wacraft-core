package websocket_model

import (
	"fmt"

	"github.com/google/uuid"
)

type ClientID struct {
	UserID uuid.UUID
	ConnID int
}

func (c *ClientID) String() string {
	return fmt.Sprintf("%s-%d", c.UserID.String(), c.ConnID)
}

func CompareClients(client1, client2 ClientID) bool {
	return client1.UserID == client2.UserID && client1.ConnID == client2.ConnID
}
