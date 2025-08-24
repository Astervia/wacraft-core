package campaign_model

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"github.com/google/uuid"
)

type QueryPaginated struct {
	Name               string     `json:"name,omitempty"`
	MessagingProductID *uuid.UUID `json:"messaging_product_id,omitempty" query:"messaging_product_id"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type QueryMessagesPaginated struct {
	CampaignID uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	MessageID  uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type QueryMessages struct {
	CampaignID uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	MessageID  uuid.UUID `json:"message_id,omitempty" query:"message_id"`

	common_model.UnrequiredID
	database_model.DateOrder
	database_model.DateWhere
}

type QueryErrorsPaginated struct {
	CampaignID        uuid.UUID `json:"campaign_id,omitempty" query:"campaign_id"`
	CampaignMessageID uuid.UUID `json:"campaign_message_id,omitempty" query:"campaign_message_id"`

	common_model.UnrequiredID
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
