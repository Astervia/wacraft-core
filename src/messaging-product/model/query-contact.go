package messaging_product_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	"github.com/google/uuid"
)

type QueryContact struct {
	ContactID          uuid.UUID `json:"contact_id,omitempty" query:"contact_id"`
	MessagingProductID uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`
	Blocked            bool      `json:"blocked,omitempty" query:"blocked"`

	common_model.UnrequiredId
	DateOrder
	DateWhere
}
