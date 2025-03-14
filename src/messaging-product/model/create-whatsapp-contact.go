package messaging_product_model

import "github.com/google/uuid"

type CreateWhatsAppContact struct {
	ProductDetails WhatsAppProductDetails `json:"product_details"`
	ContactId      uuid.UUID              `json:"contact_id"`
}
