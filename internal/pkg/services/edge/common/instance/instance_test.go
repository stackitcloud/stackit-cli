// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package instance

import (
	"fmt"
	"strings"
	"testing"

	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func TestValidateDisplayName(t *testing.T) {
	type args struct {
		displayName *string
	}
	tests := []struct {
		name string
		args *args
		want error
	}{
		// Valid cases
		{
			name: "valid minimum length",
			args: &args{displayName: utils.Ptr("test")},
		},
		{
			name: "valid maximum length",
			args: &args{displayName: utils.Ptr("testname")},
		},
		{
			name: "valid with hyphens",
			args: &args{displayName: utils.Ptr("test-app")},
		},
		{
			name: "valid with numbers",
			args: &args{displayName: utils.Ptr("test123")},
		},
		{
			name: "valid starting with letter",
			args: &args{displayName: utils.Ptr("a-test")},
		},

		// Error cases - nil pointer
		{
			name: "nil display name",
			args: &args{displayName: nil},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s may not be empty", DisplayNameFlag),
			},
		},

		// Error cases - length validation
		{
			name: "too short",
			args: &args{displayName: utils.Ptr("abc")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s is too short (minimum length is %d characters)", DisplayNameFlag, displayNameMinimumChars),
			},
		},
		{
			name: "too long",
			args: &args{displayName: utils.Ptr("verylongname")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s is too long (maximum length is %d characters)", DisplayNameFlag, displayNameMaximumChars),
			},
		},

		// Error cases - regex validation
		{
			name: "starts with number",
			args: &args{displayName: utils.Ptr("1test")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s didn't match the required regex expression %s", DisplayNameFlag, displayNameRegex),
			},
		},
		{
			name: "starts with hyphen",
			args: &args{displayName: utils.Ptr("-test")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s didn't match the required regex expression %s", DisplayNameFlag, displayNameRegex),
			},
		},
		{
			name: "ends with hyphen",
			args: &args{displayName: utils.Ptr("test-")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s didn't match the required regex expression %s", DisplayNameFlag, displayNameRegex),
			},
		},
		{
			name: "contains uppercase",
			args: &args{displayName: utils.Ptr("Test")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s didn't match the required regex expression %s", DisplayNameFlag, displayNameRegex),
			},
		},
		{
			name: "contains special characters",
			args: &args{displayName: utils.Ptr("test@")},
			want: &cliErr.FlagValidationError{
				Flag:    DisplayNameFlag,
				Details: fmt.Sprintf("%s didn't match the required regex expression %s", DisplayNameFlag, displayNameRegex),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDisplayName(tt.args.displayName)
			testUtils.AssertError(t, err, tt.want)
		})
	}
}

func TestValidatePlanId(t *testing.T) {
	type args struct {
		planId *string
	}
	tests := []struct {
		name string
		args *args
		want error
	}{
		// Valid cases
		{
			name: "valid UUID v4",
			args: &args{planId: utils.Ptr("550e8400-e29b-41d4-a716-446655440000")},
		},
		{
			name: "valid UUID lowercase",
			args: &args{planId: utils.Ptr("6ba7b810-9dad-11d1-80b4-00c04fd430c8")},
		},
		{
			name: "valid UUID uppercase",
			args: &args{planId: utils.Ptr("6BA7B810-9DAD-11D1-80B4-00C04FD430C8")},
		},
		{
			name: "valid UUID without hyphens",
			args: &args{planId: utils.Ptr("550e8400e29b41d4a716446655440000")},
		},

		// Error cases - nil pointer
		{
			name: "nil plan id",
			args: &args{planId: nil},
			want: &cliErr.FlagValidationError{
				Flag:    PlanIdFlag,
				Details: fmt.Sprintf("%s may not be empty", PlanIdFlag),
			},
		},

		// Error cases - invalid UUID format
		{
			name: "invalid UUID - too short",
			args: &args{planId: utils.Ptr("550e8400-e29b-41d4-a716")},
			want: &cliErr.FlagValidationError{
				Flag:    PlanIdFlag,
				Details: fmt.Sprintf("%s is not a valid UUID: parse 550e8400-e29b-41d4-a716 as UUID: invalid UUID length: 23", PlanIdFlag),
			},
		},
		{
			name: "invalid UUID - invalid characters",
			args: &args{planId: utils.Ptr("550e8400-e29b-41d4-a716-44665544000g")},
			want: &cliErr.FlagValidationError{
				Flag:    PlanIdFlag,
				Details: fmt.Sprintf("%s is not a valid UUID: parse 550e8400-e29b-41d4-a716-44665544000g as UUID: invalid UUID format", PlanIdFlag),
			},
		},
		{
			name: "not a UUID",
			args: &args{planId: utils.Ptr("not-a-uuid")},
			want: &cliErr.FlagValidationError{
				Flag:    PlanIdFlag,
				Details: fmt.Sprintf("%s is not a valid UUID: parse not-a-uuid as UUID: invalid UUID length: 10", PlanIdFlag),
			},
		},
		{
			name: "empty string",
			args: &args{planId: utils.Ptr("")},
			want: &cliErr.FlagValidationError{
				Flag:    PlanIdFlag,
				Details: fmt.Sprintf("%s is not a valid UUID: parse  as UUID: invalid UUID length: 0", PlanIdFlag),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlanId(tt.args.planId)
			testUtils.AssertError(t, err, tt.want)
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name string
		args *args
		want error
	}{
		// Valid cases
		{
			name: "empty description",
			args: &args{description: ""},
		},
		{
			name: "short description",
			args: &args{description: "A short description"},
		},
		{
			name: "description at maximum length",
			args: &args{description: strings.Repeat("a", descriptionMaxLength)},
		},
		{
			name: "description with special characters",
			args: &args{description: "Description with special chars: !@#$%^&*()"},
		},
		{
			name: "description with unicode",
			args: &args{description: "Description with unicode: 你好世界 🌍"},
		},

		// Error cases
		{
			name: "description too long",
			args: &args{description: strings.Repeat("a", descriptionMaxLength+1)},
			want: &cliErr.FlagValidationError{
				Flag:    DescriptionFlag,
				Details: fmt.Sprintf("%s is too long (maximum length is %d characters)", DescriptionFlag, descriptionMaxLength),
			},
		},
		{
			name: "description way too long",
			args: &args{description: strings.Repeat("a", descriptionMaxLength+100)},
			want: &cliErr.FlagValidationError{
				Flag:    DescriptionFlag,
				Details: fmt.Sprintf("%s is too long (maximum length is %d characters)", DescriptionFlag, descriptionMaxLength),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.args.description)
			testUtils.AssertError(t, err, tt.want)
		})
	}
}

func TestValidateInstanceId(t *testing.T) {
	type args struct {
		instanceId *string
	}
	tests := []struct {
		name string
		args *args
		want error
	}{
		// Valid cases
		{
			name: "valid instance id at minimum length",
			args: &args{instanceId: utils.Ptr(strings.Repeat("a", instanceIdMinLength))},
		},
		{
			name: "valid instance id at maximum length",
			args: &args{instanceId: utils.Ptr(strings.Repeat("a", instanceIdMaxLength))},
		},
		{
			name: "valid instance id with mixed characters",
			args: &args{instanceId: utils.Ptr("test-instance")},
		},

		// Error cases - nil pointer
		{
			name: "nil instance id",
			args: &args{instanceId: nil},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s may not be empty", InstanceIdFlag),
			},
		},

		// Error cases - empty string
		{
			name: "empty string",
			args: &args{instanceId: utils.Ptr("")},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s may not be empty", InstanceIdFlag),
			},
		},

		// Error cases - length validation
		{
			name: "too short",
			args: &args{instanceId: utils.Ptr(strings.Repeat("a", instanceIdMinLength-1))},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s is too short (minimum length is %d characters)", InstanceIdFlag, instanceIdMinLength),
			},
		},
		{
			name: "way too short",
			args: &args{instanceId: utils.Ptr("a")},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s is too short (minimum length is %d characters)", InstanceIdFlag, instanceIdMinLength),
			},
		},
		{
			name: "too long",
			args: &args{instanceId: utils.Ptr(strings.Repeat("a", instanceIdMaxLength+1))},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s is too long (maximum length is %d characters)", InstanceIdFlag, instanceIdMaxLength),
			},
		},
		{
			name: "way too long",
			args: &args{instanceId: utils.Ptr(strings.Repeat("a", instanceIdMaxLength+10))},
			want: &cliErr.FlagValidationError{
				Flag:    InstanceIdFlag,
				Details: fmt.Sprintf("%s is too long (maximum length is %d characters)", InstanceIdFlag, instanceIdMaxLength),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInstanceId(tt.args.instanceId)
			testUtils.AssertError(t, err, tt.want)
		})
	}
}
