package database_model

import "gorm.io/gorm"

// Provides pagination parameters.
type Paginate struct {
	Limit  int `json:"limit" default:"10" query:"limit"` // Number of items to return
	Offset int `json:"offset"default:"0" query:"offset"` // The offset from where to start the items
}

func (p *Paginate) PaginateQuery(db **gorm.DB) error {
	*db = (*db).Limit(p.Limit).Offset(p.Offset)
	return nil
}
