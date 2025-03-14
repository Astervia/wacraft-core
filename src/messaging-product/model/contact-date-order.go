package messaging_product_model

import (
	database_model "github.com/Astervia/omni-core/src/database/model"
	"gorm.io/gorm"
)

type DateOrder struct {
	LastReadAt *database_model.DateOrderEnum `json:"last_read_at,omitempty" default:"desc" query:"last_read_at"`

	database_model.DateOrder
}

func (d *DateOrder) OrderQuery(db **gorm.DB, prefix string) error {
	prefix = database_model.AddDotIfNotEmpty(prefix)

	if d.LastReadAt != nil {
		*db = (*db).Order(prefix + "last_read_at " + string(*d.LastReadAt))
	}
	d.DateOrder.OrderQuery(db, prefix)
	return nil
}
