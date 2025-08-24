package messaging_product_model

import "github.com/google/uuid"

type CreateContact struct {
	ProductDetails     ProductDetails `json:"product_details,omitempty"`
	ContactID          uuid.UUID      `json:"contact_id,omitempty"`
	MessagingProductID uuid.UUID      `json:"messaging_product_id,omitempty"`
}
