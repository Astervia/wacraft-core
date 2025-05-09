package status_entity

import (
	message_entity "github.com/Astervia/wacraft-core/src/message/entity"
	status_model "github.com/Astervia/wacraft-core/src/status/model"
)

type Status struct {
	Message *message_entity.Message `json:"json,omitempty" gorm:"foreignKey:MessageId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	status_model.StatusFields
}

func (s *Status) TableName() string {
	return "statuses"
}
