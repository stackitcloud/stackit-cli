package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Ptr Returns the pointer to any type T
func Ptr[T any](v T) *T {
	return &v
}

// CmdHelp is used to explicitly set the Run function for non-leaf commands to the command help function, so that we can catch invalid commands
// This is a workaround needed due to the open issue on the Cobra repo: https://github.com/spf13/cobra/issues/706
func CmdHelp(cmd *cobra.Command, _ []string) {
	cmd.Help() //nolint:errcheck //the function doesnt return anything to satisfy the required interface of the Run function
}

// ValidateUUID validates if the provided string is a valid UUID
func ValidateUUID(value string) error {
	_, err := uuid.Parse(value)
	if err != nil {
		return fmt.Errorf("parse %s as UUID: %w", value, err)
	}
	return nil
}

// ConvertModelToString converts an input model to a user-friendly string representation.
// This function converts the input model to a map, removes empty values, and generates a string representation of the map.
// The purpose of this function is to provide a more readable output than the default JSON representation.
// It is particularly useful when outputting to the slog logger, as the JSON format with escaped quotes does not look good.
func ConvertModelToString(model interface{}) (string, error) {
	// Marshalling and Unmarshalling is the best way to convert the struct to a map
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return "", fmt.Errorf("Error marshaling model to JSON: %v", err)
	}

	var inputModelMap map[string]interface{}
	if err := json.Unmarshal(modelBytes, &inputModelMap); err != nil {
		return "", fmt.Errorf("Error unmarshaling JSON to map: %v", err)
	}

	for key, value := range inputModelMap {
		if isEmpty(value) {
			delete(inputModelMap, key)
		}
	}

	// Generate string representation of the map
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
