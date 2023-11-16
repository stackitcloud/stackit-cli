package flags

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Returns the flag's value as a string.
// Returns "" if the flag is not set, if its value can not be converted to string, or if the flag does not exist.
func FlagToStringValue(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return ""
	}
	if cmd.Flag(flag).Changed {
		return value
	}
	return ""
}

// Returns the flag's value as a bool.
// Returns false if its value can not be converted to bool, or if the flag does not exist.
func FlagToBoolValue(cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		return false
	}
	return value
}

// Returns the flag's value as a []string.
// Returns nil if the flag is not set, if its value can not be converted to []string, or if the flag does not exist.
func FlagToStringSliceValue(cmd *cobra.Command, flag string) []string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to map[string]string, or if the flag does not exist.
func FlagToStringToStringPointer(cmd *cobra.Command, flag string) *map[string]string { //nolint:gocritic //convenient for setting the SDK payload
	value, err := cmd.Flags().GetStringToString(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to int64, or if the flag does not exist.
func FlagToInt64Pointer(cmd *cobra.Command, flag string) *int64 {
	value, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to string, or if the flag does not exist.
func FlagToStringPointer(cmd *cobra.Command, flag string) *string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to []string, or if the flag does not exist.
func FlagToStringSlicePointer(cmd *cobra.Command, flag string) *[]string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to bool, or if the flag does not exist.
func FlagToBoolPointer(cmd *cobra.Command, flag string) *bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, or if the flag does not exist.
// Returns an error if its value can not be converted to a date time with the provided format.
func FlagToDateTimePointer(cmd *cobra.Command, flag, format string) (*time.Time, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return nil, nil
	}

	if cmd.Flag(flag).Changed {
		dateTimeValue, err := time.Parse(format, value)
		if err != nil {
			return nil, fmt.Errorf("could not convert to date-time with the format %s", format)
		}
		return &dateTimeValue, nil
	}
	return nil, nil
}

// Returns the int64 value set on the flag. If no value is set, returns the flag's default value.
// Returns 0 if the flag value can not be converted to int64 or if the flag does not exist.
func FlagWithDefaultToInt64Value(cmd *cobra.Command, flag string) int64 {
	value, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		return 0
	}
	return value
}

// Returns the string value set on the flag. If no value is set, returns the flag's default value.
// Returns nil if the flag value can not be converted to string or if the flag does not exist.
func FlagWithDefaultToStringValue(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return ""
	}
	return value
}

// Returns a pointer to the flag's value. If no value is set, returns the flag's default value.
// Returns nil if the flag value can't be converted to []string or if the flag does not exist.
func FlagWithDefaultToStringSlicePointer(cmd *cobra.Command, flag string) *[]string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return nil
	}
	return &value
}
