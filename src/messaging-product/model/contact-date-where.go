package messaging_product_model

import (
	"time"

	database_model "github.com/Astervia/wacraft-core/src/database/model"
	"gorm.io/gorm"
)

type DateWhere struct {
	LastReadAtLeq *time.Time `json:"last_read_at_leq,omitempty,omitzero" query:"last_read_at_leq,omitempty,omitzero"`
	LastReadAtGeq *time.Time `json:"last_read_at_geq,omitempty,omitzero" query:"last_read_at_geq,omitempty,omitzero"`

	database_model.DateWhere
}

func (date *DateWhere) Where(db **gorm.DB, prefix string) error {
	prefix = database_model.AddDotIfNotEmpty(prefix)

	if date.LastReadAtLeq != nil {
		*db = (*db).Where(prefix+"last_read_at_leq <= ?", date.LastReadAtLeq)
	}
	if date.LastReadAtGeq != nil {
		*db = (*db).Where(prefix+"last_read_at_leq >= ?", date.LastReadAtGeq)
	}
	return date.DateWhere.Where(db, prefix)
}
