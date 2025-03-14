package message_model

import (
	message_model "github.com/Rfluid/whatsapp-cloud-api/src/message/model"
	"github.com/google/uuid"
)

type SendWhatsAppMessage struct {
	ToId       uuid.UUID             `json:"to_id"`       // Messaging product contact id to send message.
	SenderData message_model.Message `json:"sender_data"` // Specific data that allows to send message.
}
