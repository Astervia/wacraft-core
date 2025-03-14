package common_service

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func ReplacePlaceholders[U any, V any](
	toBeReplaced U,
	replaceWith V,
	name string,
) (U, error) {
	// Marshal replaceWith to JSON and unmarshal to a map
	replaceData, err := json.Marshal(replaceWith)
	if err != nil {
		var empty U
		return empty, err
	}
	var mapToReplace map[string]interface{}
	if err := json.Unmarshal(replaceData, &mapToReplace); err != nil {
		var empty U
		return empty, err
	}

	// Marshal struct1 to JSON and unmarshal to an interface{}
	dataToBeReplaced, err := json.Marshal(toBeReplaced)
	if err != nil {
		var empty U
		return empty, err
	}
	var mapToBeReplaced interface{}
	if err := json.Unmarshal(dataToBeReplaced, &mapToBeReplaced); err != nil {
		var empty U
		return empty, err
	}

	// Replace placeholders in map1 using map2
	replacePlaceholdersInMap(&mapToBeReplaced, mapToReplace, name)

	// Marshal modified map1 back to JSON
	replacedData, err := json.Marshal(mapToBeReplaced)
	if err != nil {
		var empty U
		return empty, err
	}

	// Unmarshal back to the original struct1 type
	var replacedStruct U
	if err := json.Unmarshal(replacedData, &replacedStruct); err != nil {
		var empty U
		return empty, err
	}

	// Return the modified struct1
	return replacedStruct, nil
}

func replacePlaceholdersInMap(
	data *interface{},
	mapToReplace map[string]interface{},
	name string,
) {
	switch v := (*data).(type) {
	case map[string]interface{}:
		for key, value := range v {
			replacePlaceholdersInMap(&value, mapToReplace, name)
			v[key] = value
		}
	case []interface{}:
		for i, value := range v {
			replacePlaceholdersInMap(&value, mapToReplace, name)
			v[i] = value
		}
	case string:
		// Regular expression to find placeholders like ${name.key}
		re := regexp.MustCompile(`\$\{` + regexp.QuoteMeta(name) + `\.([^}]+)\}`)
		newStr := re.ReplaceAllStringFunc(v, func(s string) string {
			// Extract the key from the placeholder
			submatches := re.FindStringSubmatch(s)
			if len(submatches) == 2 {
				key := submatches[1]
				if val, ok := mapToReplace[key]; ok {
					// Convert the value to a string
					return fmt.Sprintf("%v", val)
				}
			}
			return s // Return original string if no match
		})
		*data = newStr
	}
}
