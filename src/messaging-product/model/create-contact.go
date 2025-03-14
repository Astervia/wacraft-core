package messaging_product_model

import "github.com/google/uuid"

type CreateContact struct {
	ProductDetails     ProductDetails `json:"product_details,omitempty"`
	ContactId          uuid.UUID      `json:"contact_id,omitempty"`
	MessagingProductId uuid.UUID      `json:"messaging_product_id,omitempty"`
}
