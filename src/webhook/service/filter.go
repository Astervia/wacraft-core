package webhook_service

import (
	"encoding/json"
	"regexp"
	"strings"

	webhook_model "github.com/Astervia/wacraft-core/src/webhook/model"
)

// EvaluateFilter evaluates an event filter against a payload
// Returns true if the payload matches the filter, false otherwise
// If the filter is nil or has no conditions, returns true (allow all)
func EvaluateFilter(filter *webhook_model.EventFilter, payload any) bool {
	if filter == nil || len(filter.Conditions) == 0 {
		return true
	}

	// Convert payload to map for JSON path lookup
	payloadMap, ok := toMap(payload)
	if !ok {
		return false
	}

	// Default logic is AND
	logic := filter.Logic
	if logic == "" {
		logic = webhook_model.FilterLogicAnd
	}

	// Evaluate all conditions
	for _, condition := range filter.Conditions {
		result := evaluateCondition(condition, payloadMap)

		if logic == webhook_model.FilterLogicOr && result {
			// OR: return true on first match
			return true
		}
		if logic == webhook_model.FilterLogicAnd && !result {
			// AND: return false on first non-match
			return false
		}
	}

	// If OR and we got here, no conditions matched
	if logic == webhook_model.FilterLogicOr {
		return false
	}

	// If AND and we got here, all conditions matched
	return true
}

// evaluateCondition evaluates a single filter condition
func evaluateCondition(condition webhook_model.FilterCondition, payload map[string]any) bool {
	// Get the value at the JSON path
	value, exists := getValueByPath(payload, condition.Path)

	switch condition.Operator {
	case webhook_model.FilterOpExists:
		return exists

	case webhook_model.FilterOpEquals:
		if !exists {
			return false
		}
		return compareValues(value, condition.Value)

	case webhook_model.FilterOpContains:
		if !exists {
			return false
		}
		return containsValue(value, condition.Value)

	case webhook_model.FilterOpRegex:
		if !exists {
			return false
		}
		return matchesRegex(value, condition.Value)

	default:
		return false
	}
}

// getValueByPath retrieves a value from a nested map using dot notation
// e.g., "data.message.type" -> payload["data"]["message"]["type"]
func getValueByPath(payload map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := any(payload)

	for _, part := range parts {
		if part == "" {
			continue
		}

		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		val, exists := m[part]
		if !exists {
			return nil, false
		}
		current = val
	}

	return current, true
}

// compareValues compares two values for equality
func compareValues(a, b any) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Convert both to JSON and compare
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return string(aJSON) == string(bJSON)
}

// containsValue checks if a contains b (string containment)
func containsValue(a, b any) bool {
	aStr, aOk := toString(a)
	bStr, bOk := toString(b)

	if !aOk || !bOk {
		return false
	}

	return strings.Contains(aStr, bStr)
}

// matchesRegex checks if a value matches a regex pattern
func matchesRegex(value, pattern any) bool {
	valueStr, ok := toString(value)
	if !ok {
		return false
	}

	patternStr, ok := toString(pattern)
	if !ok {
		return false
	}

	re, err := regexp.Compile(patternStr)
	if err != nil {
		return false
	}

	return re.MatchString(valueStr)
}

// toString converts a value to string
func toString(v any) (string, bool) {
	switch val := v.(type) {
	case string:
		return val, true
	case float64:
		return strings.TrimRight(strings.TrimRight(string(rune(int(val))), "0"), "."), true
	case int:
		return string(rune(val)), true
	case bool:
		if val {
			return "true", true
		}
		return "false", true
	default:
		// Try JSON marshaling
		b, err := json.Marshal(v)
		if err != nil {
			return "", false
		}
		return string(b), true
	}
}

// toMap converts a value to map[string]any
func toMap(v any) (map[string]any, bool) {
	// If already a map, return it
	if m, ok := v.(map[string]any); ok {
		return m, true
	}

	// Try to marshal and unmarshal
	b, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, false
	}

	return m, true
}
