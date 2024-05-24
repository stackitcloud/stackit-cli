package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	pkgErrors "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
)

var cmd *cobra.Command
var service *cobra.Command
var resource *cobra.Command
var operation *cobra.Command

func setupCmd() {
	cmd = &cobra.Command{
		Use: "stackit",
	}
	service = &cobra.Command{
		Use: "service",
	}
	resource = &cobra.Command{
		Use: "resource",
	}
	operation = &cobra.Command{
		Use: "operation",
	}
	cmd.AddCommand(service)
	service.AddCommand(resource)
	resource.AddCommand(operation)
}

func TestBeautifyUnknownAndMissingCommandsError(t *testing.T) {
	tests := []struct {
		description           string
		inputError            error
		command               *cobra.Command
		expectedMsg           string
		isNotUnknownFlagError bool
	}{
		{
			description: "root command, extra input is a flag",
			inputError:  errors.New("unknown flag: --something"),
			command:     cmd,
			expectedMsg: pkgErrors.SUBCOMMAND_MISSING,
		},
		{
			description:           "non unknown flag error, return the same",
			inputError:            errors.New("some error"),
			command:               cmd,
			expectedMsg:           "some error",
			isNotUnknownFlagError: true,
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actualError := beautifyUnknownAndMissingCommandsError(cmd, tt.inputError)

			if tt.isNotUnknownFlagError {
				if actualError.Error() != tt.expectedMsg {
					t.Fatalf("expected error message to be %s, got %s", tt.expectedMsg, actualError.Error())
				}
				return
			}

			appendedErr := pkgErrors.AppendUsageTip(errors.New(tt.expectedMsg), cmd)
			if actualError.Error() != appendedErr.Error() {
				t.Fatalf("expected error to be %s, got %s", appendedErr.Error(), actualError.Error())
			}
		})
	}
}
