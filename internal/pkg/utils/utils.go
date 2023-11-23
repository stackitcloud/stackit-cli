package utils

import (
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
// Returns "false" if its value can not be converted to bool, or if the flag does not exist.
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

// Returns the int64 value set on the flag. If no value is set, returns the flag's default value.
// An error is returned if the flag value can not be converted to int64 or if the flag does not exist.
func FlagWithDefaultToInt64Value(cmd *cobra.Command, flag string) (int64, error) {
	return cmd.Flags().GetInt64(flag)
}

// Marks all given flags as required, causing the command to report an error if invoked without them.
func MarkFlagsRequired(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		err := cmd.MarkFlagRequired(flag)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ptr Returns the pointer to any type T
func Ptr[T any](v T) *T {
	return &v
}
