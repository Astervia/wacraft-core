package billing_model

import (
	database_model "github.com/Astervia/wacraft-core/src/database/model"
)

type PlanQueryPaginated struct {
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type SubscriptionQueryPaginated struct {
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type UsageLogQueryPaginated struct {
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type EndpointWeightQueryPaginated struct {
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}

type PlanPriceQueryPaginated struct {
	database_model.Paginate
	database_model.DateOrder
	database_model.DateWhere
}
