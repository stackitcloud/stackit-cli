package print

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BuildDebugStrFromInputModel converts an input model to a user-friendly string representation.
// This function converts the input model to a map, removes empty values, and generates a string representation of the map.
// The purpose of this function is to provide a more readable output than the default JSON representation.
// It is particularly useful when outputting to the slog logger, as the JSON format with escaped quotes does not look good.
func BuildDebugStrFromInputModel(model interface{}) (string, error) {
	// Marshaling and Unmarshaling is the best way to convert the struct to a map
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return "", fmt.Errorf("Error marshaling model to JSON: %w", err)
	}

	var inputModelMap map[string]interface{}
	if err := json.Unmarshal(modelBytes, &inputModelMap); err != nil {
		return "", fmt.Errorf("Error unmarshaling JSON to map: %w", err)
	}

	// Remove empty values from the map
	for key, value := range inputModelMap {
		if isEmpty(value) {
			delete(inputModelMap, key)
		}
	}

	// Build the string representation of the map
	var builder strings.Builder
	builder.WriteString("[")
	first := true
	for key, value := range inputModelMap {
		if !first {
			builder.WriteString(", ")
		} else {
			first = false
		}
		builder.WriteString(fmt.Sprintf("%s: %v", key, value))
	}
	builder.WriteString("]")
	return builder.String(), nil
}

// BuildDebugStrFromMap converts a map to a user-friendly string representation.
// This function removes empty values and generates a string representation of the map.
func BuildDebugStrFromMap(inputMap map[string]string) string {
	var builder strings.Builder
	builder.WriteString("[")
	first := true
	for key, value := range inputMap {
		if !first {
			builder.WriteString(", ")
		} else {
			first = false
		}
		builder.WriteString(fmt.Sprintf("%s: %v", key, value))
	}
	builder.WriteString("]")
	return builder.String()
}

// BuildDebugStrFromSlice converts a slice to a user-friendly string representation.
func BuildDebugStrFromSlice(inputSlice []string) string {
	sliceStr := strings.Join(inputSlice, ", ")
	return fmt.Sprintf("[%s]", sliceStr)
}

// isEmpty checks if a value is empty (nil, empty string, zero value for other types)
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}
