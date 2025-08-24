package messaging_product_model

import "gorm.io/gorm"

type ProductDetails struct {
	*WhatsAppProductDetails
}

func (p *ProductDetails) ParseIndividualFieldQueries(db **gorm.DB) {
	if p.WhatsAppProductDetails != nil {
		if p.PhoneNumber != "" {
			*db = (*db).Where("product_details->>? = ?", "phone_number", p.PhoneNumber)
		}
		if p.WaID != "" {
			*db = (*db).Where("product_details->>? = ?", "wa_id", p.WaID)
		}
	}
}
