package campaign_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	Name               string     `json:"name,omitempty"`
	MessagingProductId *uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type QueryMessagesPaginated struct {
	CampaignId uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	MessageId  uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type QueryMessages struct {
	CampaignId uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	MessageId  uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredId
	database_model.DateOrder
	database_model.DateWhere
}

type QueryErrorsPaginated struct {
	CampaignId        uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	CampaignMessageId uuid.UUID `json:"campaign_message_id,omitempty" query:"campaign_message_id"`

	common_model.UnrequiredId
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
