package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
)

var cmd *cobra.Command
var service *cobra.Command
var resource *cobra.Command
var operation *cobra.Command

var (
	testErrorMessage = "test error message"
	errStringErrTest = errors.New(testErrorMessage)
	errOpenApi404    = &oapierror.GenericOpenAPIError{StatusCode: 404, Body: []byte(`{"message":"not found"}`)}
	errOpenApi500    = &oapierror.GenericOpenAPIError{StatusCode: 500, Body: []byte(`invalid-json`)}
)

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

func setupBetaCmd() {
	cmd = &cobra.Command{
		Use: "stackit",
	}
	beta := &cobra.Command{
		Use: "beta",
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
	cmd.AddCommand(beta)
	beta.AddCommand(service)
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

func TestObservabilityInputPlanError(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		expectedMsg string
	}{
		{
			description: "base",
			args:        []string{"arg1", "arg2"},
			expectedMsg: fmt.Sprintf(OBSERVABILITY_INVALID_INPUT_PLAN, "stackit service resource operation arg1 arg2", "service"),
		},
		{
			description: "no args",
			args:        []string{},
			expectedMsg: fmt.Sprintf(OBSERVABILITY_INVALID_INPUT_PLAN, "stackit service resource operation", "service"),
		},
	}

	setupCmd()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &ObservabilityInputPlanError{
				Cmd:  operation,
				Args: tt.args,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestSetInexistentProfile(t *testing.T) {
	tests := []struct {
		description string
		profile     string
		expectedMsg string
	}{
		{
			description: "base",
			profile:     "profile",
			expectedMsg: fmt.Sprintf(SET_INEXISTENT_PROFILE, "profile", "profile"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &SetInexistentProfile{
				Profile: tt.profile,
			}

			if err.Error() != tt.expectedMsg {
				t.Fatalf("expected error to be %s, got %s", tt.expectedMsg, err.Error())
			}
		})
	}
}

func TestObservabilityInvalidPlanError(t *testing.T) {
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
			expectedMsg: fmt.Sprintf(OBSERVABILITY_INVALID_PLAN, "details", "service"),
		},
		{
			description: "no details",
			details:     "",
			service:     "service",
			expectedMsg: fmt.Sprintf(OBSERVABILITY_INVALID_PLAN, "", "service"),
		},
		{
			description: "no service",
			details:     "details",
			service:     "",
			expectedMsg: fmt.Sprintf(OBSERVABILITY_INVALID_PLAN, "details", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := &ObservabilityInvalidPlanError{
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
		service     string
		expectedMsg string
		isBetaCmd   bool
	}{
		{
			description: "no service",
			args:        []string{"arg1", "arg2"},
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_INPUT_FLAVOR, "stackit service resource operation arg1 arg2", "service"),
		},
		{
			description: "with service",
			args:        []string{"arg1", "arg2"},
			service:     "beta service",
			expectedMsg: fmt.Sprintf(DATABASE_INVALID_INPUT_FLAVOR, "stackit beta service resource operation arg1 arg2", "beta service"),
			isBetaCmd:   true,
		},
	}

	for _, tt := range tests {
		if tt.isBetaCmd {
			setupBetaCmd()
		} else {
			setupCmd()
		}
		t.Run(tt.description, func(t *testing.T) {
			err := &DatabaseInputFlavorError{
				Cmd:     operation,
				Args:    tt.args,
				Service: tt.service,
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

func TestInvalidFormatError(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				format: "",
			},
			want: "unsupported format provided",
		},
		{
			name: "with format",
			args: args{
				format: "yaml",
			},
			want: "unsupported format provided: yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&InvalidFormatError{Format: tt.args.format}).Error()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildRequestError(t *testing.T) {
	type args struct {
		reason string
		err    error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				reason: "",
				err:    nil,
			},
			want: "could not build request",
		},
		{
			name: "reason only",
			args: args{
				reason: testErrorMessage,
				err:    nil,
			},
			want: fmt.Sprintf("could not build request: %s", testErrorMessage),
		},
		{
			name: "error only",
			args: args{
				reason: "",
				err:    errStringErrTest,
			},
			want: fmt.Sprintf("could not build request: %s", testErrorMessage),
		},
		{
			name: "reason and error",
			args: args{
				reason: testErrorMessage,
				err:    errStringErrTest,
			},
			want: fmt.Sprintf("could not build request (%s): %s", testErrorMessage, testErrorMessage),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&BuildRequestError{Reason: tt.args.reason, Err: tt.args.err}).Error()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequestFailedError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil underlying",
			args: args{
				err: nil,
			},
			want: "request failed",
		},
		{
			name: "non-openapi error",
			args: args{
				err: errStringErrTest,
			},
			want: fmt.Sprintf("request failed: %s", testErrorMessage),
		},
		{
			name: "openapi error with message",
			args: args{
				err: errOpenApi404,
			},
			want: "request failed (404): not found",
		},
		{
			name: "openapi error without message",
			args: args{
				err: errOpenApi500,
			},
			want: "request failed (500): invalid-json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&RequestFailedError{Err: tt.args.err}).Error()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractMessageFromBody(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty body",
			args: args{
				body: []byte(""),
			},
			want: "",
		},
		{
			name: "invalid json",
			args: args{
				body: []byte("not-json"),
			},
			want: "",
		},
		{
			name: "missing message field",
			args: args{
				body: []byte(`{"error":"oops"}`),
			},
			want: "",
		},
		{
			name: "with message field",
			args: args{
				body: []byte(`{"message":"the reason"}`),
			},
			want: "the reason",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOpenApiMessageFromBody(tt.args.body)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConstructorsReturnExpected(t *testing.T) {
	buildRequestError := NewBuildRequestError(testErrorMessage, errStringErrTest)

	tests := []struct {
		name string
		got  any
		want any
	}{
		{
			name: "InvalidFormat format",
			got:  NewInvalidFormatError("fmt").Format,
			want: "fmt",
		},
		{
			name: "BuildRequestError error",
			got:  buildRequestError.Err,
			want: errStringErrTest,
		},
		{
			name: "BuildRequestError reason",
			got:  buildRequestError.Reason,
			want: testErrorMessage,
		},
		{
			name: "RequestFailed error",
			got:  NewRequestFailedError(errStringErrTest).Err,
			want: errStringErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantErr, wantIsErr := tt.want.(error)
			gotErr, gotIsErr := tt.got.(error)
			if wantIsErr {
				if !gotIsErr {
					t.Fatalf("expected error but got %T", tt.got)
				}
				if !errors.Is(gotErr, wantErr) {
					t.Errorf("got error %v, want %v", gotErr, wantErr)
				}
				return
			}

			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
