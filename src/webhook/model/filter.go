package webhook_model

// FilterOperator represents the comparison operator for event filtering
type FilterOperator string

const (
	FilterOpEquals   FilterOperator = "equals"
	FilterOpContains FilterOperator = "contains"
	FilterOpRegex    FilterOperator = "regex"
	FilterOpExists   FilterOperator = "exists"
)

// FilterLogic represents the logical operator for combining conditions
type FilterLogic string

const (
	FilterLogicAnd FilterLogic = "AND"
	FilterLogicOr  FilterLogic = "OR"
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Path     string         `json:"path"`            // JSON path to the field (e.g., "data.type")
	Operator FilterOperator `json:"operator"`        // Comparison operator
	Value    any            `json:"value,omitempty"` // Value to compare against (not needed for "exists")
}

// EventFilter represents the filter configuration for webhook events
type EventFilter struct {
	Logic      FilterLogic       `json:"logic,omitempty"`      // AND or OR (default: AND)
	Conditions []FilterCondition `json:"conditions,omitempty"` // List of conditions to evaluate
}
