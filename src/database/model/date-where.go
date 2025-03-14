package database_model

import (
	"time"

	"gorm.io/gorm"
)

type DateWhere struct {
	CreatedAdLeq time.Time `json:"created_at_leq,omitempty,omitzero" query:"created_at_leq,omitempty,omitzero"`
	UpdatedAtLeq time.Time `json:"updated_at_leq,omitempty,omitzero" query:"updated_at_leq,omitempty,omitzero"`
	CreatedAtGeq time.Time `json:"created_at_geq,omitempty,omitzero" query:"created_at_geq,omitempty,omitzero"`
	UpdatedAtGeq time.Time `json:"updated_at_geq,omitempty,omitzero" query:"updated_at_geq,omitempty,omitzero"`
}

func (date *DateWhere) Where(db **gorm.DB, prefix string) error {
	prefix = AddDotIfNotEmpty(prefix)

	if !date.CreatedAdLeq.IsZero() {
		*db = (*db).Where(prefix+"created_at <= ?", date.CreatedAdLeq)
	}
	if !date.UpdatedAtLeq.IsZero() {
		*db = (*db).Where(prefix+"updated_at <= ?", date.UpdatedAtLeq)
	}
	if !date.CreatedAtGeq.IsZero() {
		*db = (*db).Where(prefix+"created_at >= ?", date.CreatedAtGeq)
	}
	if !date.UpdatedAtGeq.IsZero() {
		*db = (*db).Where(prefix+"updated_at >= ?", date.UpdatedAtGeq)
	}
	return nil
}

type DateWhereWithDeletedAt struct {
	DeletedAtLeq time.Time `json:"deleted_at_leq,omitempty,omitzero" query:"deleted_at_leq,omitempty,omitzero"`
	DeletedAtGeq time.Time `json:"deleted_at_geq,omitempty,omitzero" query:"deleted_at_geq,omitempty,omitzero"`
	DateWhere
}

func (date *DateWhereWithDeletedAt) Where(db **gorm.DB, prefix string) error {
	prefix = AddDotIfNotEmpty(prefix)

	if !date.DeletedAtLeq.IsZero() {
		*db = (*db).Where(prefix+"deleted_at <= ?", date.DeletedAtLeq)
	}
	if !date.DeletedAtGeq.IsZero() {
		*db = (*db).Where(prefix+"deleted_at >= ?", date.DeletedAtGeq)
	}
	return date.DateWhere.Where(db, prefix)
}

func AddDotIfNotEmpty(prefix string) string {
	if prefix != "" {
		prefix += "."
	}
	return prefix
}
