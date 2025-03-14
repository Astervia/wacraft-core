package campaign_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	"github.com/google/uuid"
)

type UpdateCampaign struct {
	Name               string     `json:"name,omitempty"`
	MessagingProductId *uuid.UUID `json:"messaging_product_id,omitempty"`

	common_model.RequiredId
}
