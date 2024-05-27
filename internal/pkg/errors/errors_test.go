package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
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

func TestSimpleErrors(t *testing.T) {
	tests := []struct {
		description string
		err         error
		expectedMsg string
	}{
		{
			description: "Test ProjectIdError",
			err:         &ProjectIdError{},
			expectedMsg: MISSING_PROJECT_ID,
		},
		{
			description: "Test EmptyUpdateError",
			err:         &EmptyUpdateError{},
			expectedMsg: EMPTY_UPDATE,
		},
		{
			description: "Test AuthError",
			err:         &AuthError{},
			expectedMsg: FAILED_AUTH,
		},
		{
			description: "Test ActivateServiceAccountError",
			err:         &ActivateServiceAccountError{},
			expectedMsg: FAILED_SERVICE_ACCOUNT_ACTIVATION,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if tt.err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, tt.err.Error())
			}
		})
	}
}

func TestArgusInputPlanError(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		expectedMsg string
	}{
		{
			description: "base",
			args:        []string{"arg1", "arg2"},
			expectedMsg: fmt.Sprintf(ARGUS_INVALID_INPUT_PLAN, "stackit service resource operation arg1 arg2", "service"),
		},
		{
			description: "no args",
			args:        []string{},
			expectedMsg: fmt.Sprintf(ARGUS_INVALID_INPUT_PLAN, "stackit service resource operation", "service"),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &ArgusInputPlanError{
				Cmd:  operation,
				Args: tt.args,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestArgusInvalidPlanError(t *testing.T) {
	tests := []struct {
		description string
		details     string
		service     string
		expectedMsg string
	}{
		{
			description: "base",
			details:     "details",
			service:     "service",
			expectedMsg: fmt.Sprintf(ARGUS_INVALID_PLAN, "details", "service"),
		},
		{
			description: "no details",
			details:     "",
			service:     "service",
			expectedMsg: fmt.Sprintf(ARGUS_INVALID_PLAN, "", "service"),
		},
		{
			description: "no service",
			details:     "details",
			service:     "",
			expectedMsg: fmt.Sprintf(ARGUS_INVALID_PLAN, "details", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &ArgusInvalidPlanError{
				Service: tt.service,
				Details: tt.details,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestDSAInputPlanError(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		expectedMsg string
	}{
		{
			description: "base",
			args:        []string{"arg1", "arg2"},
			expectedMsg: fmt.Sprintf(DSA_INVALID_INPUT_PLAN, "stackit service resource operation arg1 arg2", "service"),
		},
		{
			description: "no args",
			args:        []string{},
			expectedMsg: fmt.Sprintf(DSA_INVALID_INPUT_PLAN, "stackit service resource operation", "service"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			setupCmd()
			err := &DSAInputPlanError{
				Cmd:  operation,
				Args: tt.args,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestDSAInvalidPlanError(t *testing.T) {
	tests := []struct {
		description string
		details     string
		service     string
		expectedMsg string
	}{
		{
			description: "base",
			details:     "details",
			service:     "service",
			expectedMsg: fmt.Sprintf(DSA_INVALID_PLAN, "details", "service"),
		},
		{
			description: "no details",
			details:     "",
			service:     "service",
			expectedMsg: fmt.Sprintf(DSA_INVALID_PLAN, "", "service"),
		},
		{
			description: "no service",
			details:     "details",
			service:     "",
			expectedMsg: fmt.Sprintf(DSA_INVALID_PLAN, "details", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &DSAInvalidPlanError{
				Service: tt.service,
				Details: tt.details,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestDatabaseInputFlavorError(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		operation   string
		expectedMsg string
	}{
		{
			description: "base",
			args:        []string{"arg1", "arg2"},
			operation:   "operation",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_INPUT_FLAVOR, "stackit service resource operation arg1 arg2", "service"),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &DatabaseInputFlavorError{
				Cmd:       operation,
				Args:      tt.args,
				Operation: tt.operation,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestDatabaseInvalidFlavorError(t *testing.T) {
	tests := []struct {
		description string
		details     string
		service     string
		expectedMsg string
	}{
		{
			description: "base",
			details:     "details",
			service:     "service",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_FLAVOR, "details", "service"),
		},
		{
			description: "no details",
			details:     "",
			service:     "service",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_FLAVOR, "", "service"),
		},
		{
			description: "no service",
			details:     "details",
			service:     "",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_FLAVOR, "details", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &DatabaseInvalidFlavorError{
				Service: tt.service,
				Details: tt.details,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestDatabaseInvalidStorageError(t *testing.T) {
	tests := []struct {
		description string
		details     string
		service     string
		flavorId    string
		expectedMsg string
	}{
		{
			description: "base",
			details:     "details",
			service:     "service",
			flavorId:    "flavorId",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_STORAGE, "details", "service", "flavorId"),
		},
		{
			description: "no details",
			details:     "",
			service:     "service",
			flavorId:    "flavorId",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_STORAGE, "", "service", "flavorId"),
		},
		{
			description: "no service",
			details:     "details",
			service:     "",
			flavorId:    "flavorId",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_STORAGE, "details", "", "flavorId"),
		},
		{
			description: "no flavorId",
			details:     "details",
			service:     "service",
			flavorId:    "",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_STORAGE, "details", "service", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &DatabaseInvalidStorageError{
				Service:  tt.service,
				Details:  tt.details,
				FlavorId: tt.flavorId,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestFlagValidationError(t *testing.T) {
	tests := []struct {
		description string
		flag        string
		details     string
		expectedMsg string
	}{
		{
			description: "base",
			flag:        "flag",
			details:     "details",
			expectedMsg: fmt.Sprintf(FLAG_VALIDATION, "flag", "details"),
		},
		{
			description: "no flag",
			flag:        "",
			details:     "details",
			expectedMsg: fmt.Sprintf(FLAG_VALIDATION, "", "details"),
		},
		{
			description: "no details",
			flag:        "flag",
			details:     "",
			expectedMsg: fmt.Sprintf(FLAG_VALIDATION, "flag", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &FlagValidationError{
				Flag:    tt.flag,
				Details: tt.details,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestRequiredMutuallyExclusiveFlagsError(t *testing.T) {
	tests := []struct {
		description string
		flags       []string
		expectedMsg string
	}{
		{
			description: "base",
			flags:       []string{"flag1", "flag2"},
			expectedMsg: fmt.Sprintf(REQUIRED_MUTUALLY_EXCLUSIVE_FLAGS, "flag1, flag2"),
		},
		{
			description: "no flags",
			flags:       []string{},
			expectedMsg: fmt.Sprintf(REQUIRED_MUTUALLY_EXCLUSIVE_FLAGS, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &RequiredMutuallyExclusiveFlagsError{
				Flags: tt.flags,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestArgValidationError(t *testing.T) {
	tests := []struct {
		description string
		arg         string
		details     string
		expectedMsg string
	}{
		{
			description: "base",
			arg:         "arg",
			details:     "details",
			expectedMsg: fmt.Sprintf(ARG_VALIDATION, "arg", "details"),
		},
		{
			description: "no arg",
			arg:         "",
			details:     "details",
			expectedMsg: fmt.Sprintf(ARG_VALIDATION, "", "details"),
		},
		{
			description: "no details",
			arg:         "arg",
			details:     "",
			expectedMsg: fmt.Sprintf(ARG_VALIDATION, "arg", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &ArgValidationError{
				Arg:     tt.arg,
				Details: tt.details,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestSingleArgExpectedError(t *testing.T) {
	tests := []struct {
		description string
		expected    string
		count       int
		expectedMsg string
	}{
		{
			description: "base",
			expected:    "expected",
			count:       1,
			expectedMsg: fmt.Sprintf(ARG_MISSING, "expected"),
		},
		{
			description: "multiple",
			expected:    "expected",
			count:       2,
			expectedMsg: fmt.Sprintf(SINGLE_ARG_EXPECTED, "expected", 2),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &SingleArgExpectedError{
				Expected: tt.expected,
				Count:    tt.count,
				Cmd:      operation,
			}

			appendedErr := AppendUsageTip(errors.New(tt.expectedMsg), operation)

			if err.Error() != appendedErr.Error() {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestSingleOptionalArgExpectedError(t *testing.T) {
	tests := []struct {
		description string
		expected    string
		count       int
		expectedMsg string
	}{
		{
			description: "base",
			expected:    "expected",
			count:       1,
			expectedMsg: fmt.Sprintf(SINGLE_OPTIONAL_ARG_EXPECTED, "expected", 1),
		},
		{
			description: "multiple",
			expected:    "expected",
			count:       2,
			expectedMsg: fmt.Sprintf(SINGLE_OPTIONAL_ARG_EXPECTED, "expected", 2),
		},
		{
			description: "no count",
			expected:    "expected",
			count:       0,
			expectedMsg: fmt.Sprintf(SINGLE_OPTIONAL_ARG_EXPECTED, "expected", 0),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &SingleOptionalArgExpectedError{
				Expected: tt.expected,
				Count:    tt.count,
				Cmd:      operation,
			}

			appendedErr := AppendUsageTip(errors.New(tt.expectedMsg), operation)

			if err.Error() != appendedErr.Error() {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestInputUnknownError(t *testing.T) {
	tests := []struct {
		description string
		input       string
		command     *cobra.Command
		expectedMsg string
	}{
		{
			description: "extra argument, not a subcommand",
			input:       "extra",
			command:     operation,
			expectedMsg: fmt.Sprintf(ARG_UNKNOWN, "extra"),
		},
		{
			description: "extra subcommand",
			input:       "extra",
			command:     service,
			expectedMsg: fmt.Sprintf(SUBCOMMAND_UNKNOWN, "extra"),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &InputUnknownError{
				ProvidedInput: tt.input,
				Cmd:           tt.command,
			}

			appendedErr := AppendUsageTip(errors.New(tt.expectedMsg), tt.command)

			if err.Error() != appendedErr.Error() {
				t.Fatalf("expected error to be %s, got %s", appendedErr.Error(), err.Error())
			}
		})
	}
}

func TestSubcommandMissingError(t *testing.T) {
	tests := []struct {
		description string
		expectedMsg string
	}{
		{
			description: "base",
			expectedMsg: SUBCOMMAND_MISSING,
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &SubcommandMissingError{
				Cmd: cmd,
			}

			appendedErr := AppendUsageTip(errors.New(tt.expectedMsg), cmd)

			if err.Error() != appendedErr.Error() {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestAppendUsageTip(t *testing.T) {
	tests := []struct {
		description   string
		err           error
		expectedError error
	}{
		{
			description:   "base",
			err:           fmt.Errorf("error"),
			expectedError: fmt.Errorf("%w.\n\n%s", fmt.Errorf("error"), fmt.Sprintf(USAGE_TIP, "stackit")),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := AppendUsageTip(tt.err, cmd)

			if err.Error() != tt.expectedError.Error() {
				t.Fatalf("expected error to be %s, got %s", tt.expectedError, err.Error())
			}
		})
	}
}
