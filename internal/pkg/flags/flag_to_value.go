package flags

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

// Returns the flag's value as a string.
// Returns "" if the flag is not set, if its value can not be converted to string, or if the flag does not exist.
func FlagToStringValue(p *print.Printer, cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to string value: %v", err)
		return ""
	}
	if cmd.Flag(flag).Changed {
		return value
	}
	return ""
}

// Returns the flag's value as a bool.
// Returns false if its value can not be converted to bool, or if the flag does not exist.
func FlagToBoolValue(p *print.Printer, cmd *cobra.Command, flag string) bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to bool value: %v", err)
		return false
	}
	return value
}

// Returns the flag's value as a []string.
// Returns nil if the flag is not set, if its value can not be converted to []string, or if the flag does not exist.
func FlagToStringSliceValue(p *print.Printer, cmd *cobra.Command, flag string) []string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to string slice value: %v", err)
		return nil
	}
	if cmd.Flag(flag).Changed {
		return value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to map[string]string, or if the flag does not exist.
func FlagToStringToStringPointer(p *print.Printer, cmd *cobra.Command, flag string) *map[string]string { //nolint:gocritic //convenient for setting the SDK payload
	value, err := cmd.Flags().GetStringToString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to string to string pointer: %v", err)
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to int64, or if the flag does not exist.
func FlagToInt64Pointer(p *print.Printer, cmd *cobra.Command, flag string) *int64 {
	value, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to Int64 pointer: %v", err)
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to string, or if the flag does not exist.
func FlagToStringPointer(p *print.Printer, cmd *cobra.Command, flag string) *string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to string pointer: %v", err)
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to []string, or if the flag does not exist.
func FlagToStringSlicePointer(p *print.Printer, cmd *cobra.Command, flag string) *[]string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to string slice pointer: %v", err)
		return nil
	}
	if cmd.Flag(flag).Changed {
		return &value
	}
	return nil
}

// Returns a pointer to the flag's value.
// Returns nil if the flag is not set, if its value can not be converted to bool, or if the flag does not exist.
func FlagToBoolPointer(p *print.Printer, cmd *cobra.Command, flag string) *bool {
	value, err := cmd.Flags().GetBool(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to bool pointer: %v", err)
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
func FlagToDateTimePointer(p *print.Printer, cmd *cobra.Command, flag, format string) (*time.Time, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag to date-time pointer: %v", err)
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
func FlagWithDefaultToInt64Value(p *print.Printer, cmd *cobra.Command, flag string) int64 {
	value, err := cmd.Flags().GetInt64(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag with default to Int64 value: %v", err)
		return 0
	}
	return value
}

// Returns the string value set on the flag. If no value is set, returns the flag's default value.
// Returns nil if the flag value can not be converted to string or if the flag does not exist.
func FlagWithDefaultToStringValue(p *print.Printer, cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag with default to string value: %v", err)
		return ""
	}
	return value
}

// Returns a pointer to the flag's value. If no value is set, returns the flag's default value.
// Returns nil if the flag value can't be converted to []string or if the flag does not exist.
func FlagWithDefaultToStringSlicePointer(p *print.Printer, cmd *cobra.Command, flag string) *[]string {
	value, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		p.Debug(print.ErrorLevel, "convert flag with default to string slice pointer: %v", err)
		return nil
	}
	return &value
}
