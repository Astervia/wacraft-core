package contact_entity

import (
	common_model "github.com/Astervia/wacraft-core/src/common/model"
)

type Contact struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	PhotoPath string `json:"photo_path,omitempty"`

	common_model.Audit
}
