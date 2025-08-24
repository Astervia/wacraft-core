package messaging_product_model

type WhatsAppProductDetails struct {
	PhoneNumber string `json:"phone_number" query:"phone_number"` // Available at from field on received messages.
	WaID        string `json:"wa_id" query:"wa_id"`               // Available at from field on received messages.
}

type UnrequiredWhatsAppProductDetails struct {
	PhoneNumber string `json:"phone_number,omitempty" query:"phone_number"`
	WaID        string `json:"wa_id,omitempty" query:"wa_id"`
}
