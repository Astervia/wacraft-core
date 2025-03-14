package format_service

import (
	"encoding/json"
)

func Json[T any](
	jsonValue *T,
) (string, error) {
	jsonAsByte, err := json.MarshalIndent(jsonValue, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonAsByte), err
}
