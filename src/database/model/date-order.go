package database_model

import "gorm.io/gorm"

type DateOrderEnum string

const (
	Asc  DateOrderEnum = "asc"
	Desc DateOrderEnum = "desc"
)

type DateOrder struct {
	CreatedAt *DateOrderEnum `json:"created_at,omitempty" default:"desc" query:"created_at"`
	UpdatedAt *DateOrderEnum `json:"updated_at,omitempty" default:"desc" query:"updated_at"`
}

type DateOrderWithDeletedAt struct {
	DeletedAt *DateOrderEnum `json:"deleted_at,omitempty" default:"desc" query:"deleted_at"`
	DateOrder
}

func (d *DateOrder) OrderQuery(db **gorm.DB, prefix string) error {
	prefix = AddDotIfNotEmpty(prefix)

	if d.CreatedAt != nil {
		*db = (*db).Order(prefix + "created_at " + string(*d.CreatedAt))
	}
	if d.UpdatedAt != nil {
		*db = (*db).Order(prefix + "updated_at " + string(*d.UpdatedAt))
	}
	return nil
}

func (d *DateOrderWithDeletedAt) OrderQuery(db **gorm.DB, prefix string) error {
	prefix = AddDotIfNotEmpty(prefix)

	if d.DeletedAt != nil {
		*db = (*db).Order(prefix + "deleted_at " + string(*d.DeletedAt))
	}
	d.DateOrder.OrderQuery(db, prefix)
	return nil
}
