package utils

import (
	"fmt"

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
