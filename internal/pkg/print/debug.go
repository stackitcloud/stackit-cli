package print

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// BuildDebugStrFromInputModel converts an input model to a user-friendly string representation.
// This function converts the input model to a map, removes empty values, and generates a string representation of the map.
// The purpose of this function is to provide a more readable output than the default JSON representation.
// It is particularly useful when outputting to the slog logger, as the JSON format with escaped quotes does not look good.
func BuildDebugStrFromInputModel(model any) (string, error) {
	// Marshaling and Unmarshaling is the best way to convert the struct to a map
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return "", fmt.Errorf("marshal model to JSON: %w", err)
	}

	var inputModelMap map[string]any
	if err := json.Unmarshal(modelBytes, &inputModelMap); err != nil {
		return "", fmt.Errorf("unmarshaling JSON to map: %w", err)
	}

	return BuildDebugStrFromMap(inputModelMap), nil
}

// BuildDebugStrFromMap converts a map to a user-friendly string representation.
// This function removes empty values and generates a string representation of the map.
// The string representation is in the format: [key1: value1, key2: value2, ...]
// The keys are ordered alphabetically to make the output deterministic.
func BuildDebugStrFromMap(inputMap map[string]any) string {
	// Sort the keys to make the output deterministic
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := "["
	first := true
	for _, key := range keys {
		value := inputMap[key]
		if isEmpty(value) {
			continue
		}

		if !first {
			result += ", "
		} else {
			first = false
		}

		result += fmt.Sprintf("%s: %v", key, value)
	}
	result += "]"
	return result
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
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []bool:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	default:
		return false
	}
}
