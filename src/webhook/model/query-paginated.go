package webhook_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
	database_model "github.com/Astervia/omni-core/src/database/model"
	"github.com/google/uuid"
)

// QueryPaginated represents the paginated query structure for webhooks
type QueryPaginated struct {
	Url        string `json:"url,omitempty"`                             // Optional URL filter
	Event      Event  `json:"event,omitempty"`                           // Optional event filter
	HttpMethod string `json:"http_method,omitempty" query:"http_method"` // Optional HTTP method filter
	Timeout    *int   `json:"timeout,omitempty"`

	database_model.Paginate   // Pagination structure (e.g., limit, offset)
	database_model.DateOrder  // Date ordering options (ASC/DESC)
	database_model.DateWhere  // Date filtering (e.g., before/after specific date)
	common_model.UnrequiredId // Optional ID for querying by specific webhook ID
}

type QueryLogsPaginated struct {
	WebhookId        uuid.UUID `json:"webhook_id" query:"webhook_id"`
	HttpResponseCode int       `json:"http_response_code,omitempty" query:"http_response_code"`

	database_model.Paginate   // Pagination structure (e.g., limit, offset)
	database_model.DateOrder  // Date ordering options (ASC/DESC)
	database_model.DateWhere  // Date filtering (e.g., before/after specific date)
	common_model.UnrequiredId // Optional ID for querying by specific webhook ID
}
