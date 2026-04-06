package campaign_model

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleCampaign sets a campaign's scheduled_at time and transitions it to "scheduled" status.
type ScheduleCampaign struct {
	ID          uuid.UUID  `json:"id" validate:"required"`
	ScheduledAt *time.Time `json:"scheduled_at" validate:"required"`
}

// UnscheduleCampaign cancels a pending schedule, resetting the campaign to "draft" status.
type UnscheduleCampaign struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
