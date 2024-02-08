package args

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/spf13/cobra"
)

func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}

	return &errors.InputUnknownError{
		ProvidedInput: args[0],
		Cmd:           cmd,
	}
}

// SingleArg checks if only one non-empty argument was provided and validates it
// using the validate function. It returns an error if none or multiple arguments
// are provided, or if the argument is invalid.
// For no validation, you can pass a nil validate function
func SingleArg(argName string, validate func(value string) error) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 || args[0] == "" {
			return &errors.SingleArgExpectedError{
				Cmd:      cmd,
				Expected: argName,
				Count:    len(args),
			}
		}
		if validate != nil {
			err := validate(args[0])
			if err != nil {
				return &errors.ArgValidationError{
					Arg:     argName,
					Details: err.Error(),
				}
			}
		}
		return nil
	}
}
