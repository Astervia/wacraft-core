package billing_model

import "github.com/google/uuid"

// UsageSummary represents the current throughput usage for a scope.
type UsageSummary struct {
	Scope           Scope      `json:"scope"`
	UserID          *uuid.UUID `json:"user_id,omitempty"`
	WorkspaceID     *uuid.UUID `json:"workspace_id,omitempty"`
	Unlimited       bool       `json:"unlimited"`        // True if scope has infinite throughput
	ThroughputLimit int        `json:"throughput_limit"` // Total allowed weighted requests per window (0 when unlimited)
	WindowSeconds   int        `json:"window_seconds"`   // Window duration
	CurrentUsage    int64      `json:"current_usage"`    // Weighted requests used in current window
	Remaining       int64      `json:"remaining"`        // Requests remaining (-1 when unlimited)
	Fallback        bool       `json:"fallback"`         // True if this entry represents the fallback budget for billing routes
}
