package contact_model

import (
	common_model "github.com/Astervia/omni-core/src/common/model"
)

type CreateContact struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	PhotoPath string `json:"photo_path,omitempty"`
}

type UpdateContact struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	PhotoPath string `json:"photo_path,omitempty"`

	common_model.RequiredId
}
