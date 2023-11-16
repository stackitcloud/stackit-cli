package flags

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"stackit/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func TestEnumFlag(t *testing.T) {
	options := []string{"foo", "BaR"}

	tests := []struct {
		description string
		ignoreCase  bool
		value       string
		isValid     bool
	}{
		{
			description: "valid",
			value:       "foo",
			isValid:     true,
		},
		{
			description: "empty",
			value:       "",
			isValid:     false,
		},
		{
			description: "invalid 1",
			value:       "ba",
			isValid:     false,
		},
		{
			description: "invalid 2",
			value:       "foo ",
			isValid:     false,
		},
		{
			description: "invalid 3",
			value:       "bar",
			isValid:     false,
		},
		{
			description: "ignore case - valid 1",
			ignoreCase:  true,
			value:       "foo",
			isValid:     true,
		},
		{
			description: "ignore case - valid 2",
			ignoreCase:  true,
			value:       "fOO",
			isValid:     true,
		},
		{
			description: "ignore case - valid 3",
			ignoreCase:  true,
			value:       "bar",
			isValid:     true,
		},
		{
			description: "ignore case - invalid 1",
			ignoreCase:  true,
			value:       "ba",
			isValid:     false,
		},
		{
			description: "ignore case - invalid 2",
			ignoreCase:  true,
			value:       "foo ",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := EnumFlag(tt.ignoreCase, "", options...)
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			err := cmd.Flags().Set("test-flag", tt.value)

			if !tt.isValid && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			value := FlagToStringValue(cmd, "test-flag")
			if !tt.ignoreCase && value != tt.value {
				t.Fatalf("flag did not return set value")
			}
			if tt.ignoreCase && !strings.EqualFold(value, tt.value) {
				t.Fatalf("flag did not return set value")
			}
		})
	}
}

func TestEnumSliceFlag(t *testing.T) {
	validOption1 := "foo"
	validOption2 := "BaR"
	validOption3 := "baz"

	validOption2Lower := strings.ToLower(validOption2)

	options := []string{validOption1, validOption2, validOption3}

	tests := []struct {
		description   string
		ignoreCase    bool
		defaultValue  []string
		value1        *string
		value2        *string
		expectedValue []string
		isValid       bool
	}{
		{
			description:   "valid two single values",
			value1:        utils.Ptr(validOption1),
			value2:        utils.Ptr(validOption2),
			expectedValue: []string{validOption1, validOption2},
			isValid:       true,
		},
		{
			description:   "valid list value",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", validOption1, validOption2)),
			expectedValue: []string{validOption1, validOption2},
			isValid:       true,
		},
		{
			description:   "valid list value and single value",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", validOption1, validOption2)),
			value2:        utils.Ptr(validOption3),
			expectedValue: []string{validOption1, validOption2, validOption3},
			isValid:       true,
		},
		{
			description:   "valid two list values",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", validOption1, validOption2)),
			value2:        utils.Ptr(fmt.Sprintf("%s,%s", validOption2, validOption3)),
			expectedValue: []string{validOption1, validOption2, validOption2, validOption3},
			isValid:       true,
		},
		{
			description: "invalid value",
			value1:      utils.Ptr("invalid-value"),
			value2:      utils.Ptr(validOption1),
			isValid:     false,
		},
		{
			description: "invalid value in list",
			value1:      utils.Ptr(fmt.Sprintf("invalid-value,%s", validOption1)),
			isValid:     false,
		},
		{
			description: "invalid empty value",
			value1:      utils.Ptr(""),
			isValid:     false,
		},
		{
			description: "invalid empty value in list",
			value1:      utils.Ptr(fmt.Sprintf(",%s", validOption1)),
			isValid:     false,
		},
		{
			description:   "no values",
			expectedValue: []string{},
			isValid:       true,
		},
		{
			description:   "ignore case - valid single value",
			value1:        utils.Ptr(validOption2Lower),
			ignoreCase:    true,
			expectedValue: []string{validOption2Lower},
			isValid:       true,
		},
		{
			description:   "ignore case - valid in list",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", validOption1, validOption2Lower)),
			ignoreCase:    true,
			expectedValue: []string{validOption1, validOption2Lower},
			isValid:       true,
		},
		{
			description: "ignore case - invalid single value",
			value1:      utils.Ptr("ba"),
			ignoreCase:  true,
			isValid:     false,
		},
		{
			description: "ignore case - invalid in list",
			value1:      utils.Ptr(fmt.Sprintf("%s,%s", validOption1, "ba")),
			ignoreCase:  true,
			isValid:     false,
		},
		{
			description:   "default value",
			defaultValue:  []string{validOption1, validOption2},
			expectedValue: []string{validOption1, validOption2},
			isValid:       true,
		},
		{
			description:   "default value - set value",
			defaultValue:  []string{validOption1, validOption2},
			value1:        utils.Ptr(validOption1),
			expectedValue: []string{validOption1},
			isValid:       true,
		},
		{
			description:   "ignore case - default value",
			defaultValue:  []string{validOption2},
			ignoreCase:    true,
			expectedValue: []string{validOption2Lower},
			isValid:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := EnumSliceFlag(tt.ignoreCase, tt.defaultValue, options...)
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			var err1, err2 error
			if tt.value1 != nil {
				err1 = cmd.Flags().Set("test-flag", *tt.value1)
			}
			if tt.value2 != nil {
				err2 = cmd.Flags().Set("test-flag", *tt.value2)
			}

			if !tt.isValid && err1 == nil && err2 == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err1 != nil {
				t.Fatalf("failed on valid input: %v", err1)
			}
			if err2 != nil {
				t.Fatalf("failed on valid input: %v", err2)
			}
			value, err := cmd.Flags().GetStringSlice("test-flag")
			if err != nil {
				t.Fatalf("failed to get value: %v", err)
			}
			if !reflect.DeepEqual(tt.expectedValue, value) {
				t.Fatalf("flag did not return set value (expected %s, got %s)", tt.expectedValue, value)
			}
		})
	}
}

func TestUUIDFlag(t *testing.T) {
	tests := []struct {
		description string
		value       string
		isValid     bool
	}{
		{
			description: "valid",
			value:       uuid.NewString(),
			isValid:     true,
		},
		{
			description: "empty",
			value:       "",
			isValid:     false,
		},
		{
			description: "invalid",
			value:       "invalid-uuid",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := UUIDFlag()
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			err := cmd.Flags().Set("test-flag", tt.value)

			if !tt.isValid && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			value := FlagToStringValue(cmd, "test-flag")
			if value != tt.value {
				t.Fatalf("flag did not return set value")
			}
		})
	}
}

func TestUUIDSliceFlag(t *testing.T) {
	testUUID1 := uuid.NewString()
	testUUID2 := uuid.NewString()
	testUUID3 := uuid.NewString()
	testUUID4 := uuid.NewString()
	tests := []struct {
		description   string
		value1        *string
		value2        *string
		expectedValue []string
		isValid       bool
	}{
		{
			description:   "valid two single values",
			value1:        utils.Ptr(testUUID1),
			value2:        utils.Ptr(testUUID2),
			expectedValue: []string{testUUID1, testUUID2},
			isValid:       true,
		},
		{
			description:   "valid list value",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", testUUID1, testUUID2)),
			expectedValue: []string{testUUID1, testUUID2},
			isValid:       true,
		},
		{
			description:   "valid list value and single value",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", testUUID1, testUUID2)),
			value2:        utils.Ptr(testUUID3),
			expectedValue: []string{testUUID1, testUUID2, testUUID3},
			isValid:       true,
		},
		{
			description:   "valid two list values",
			value1:        utils.Ptr(fmt.Sprintf("%s,%s", testUUID1, testUUID2)),
			value2:        utils.Ptr(fmt.Sprintf("%s,%s", testUUID3, testUUID4)),
			expectedValue: []string{testUUID1, testUUID2, testUUID3, testUUID4},
			isValid:       true,
		},
		{
			description: "invalid value",
			value1:      utils.Ptr("invalid-UUID"),
			value2:      utils.Ptr(testUUID1),
			isValid:     false,
		},
		{
			description: "invalid value in list",
			value1:      utils.Ptr(fmt.Sprintf("invalid-UUID,%s", testUUID1)),
			isValid:     false,
		},
		{
			description: "invalid empty value",
			value1:      utils.Ptr(""),
			isValid:     false,
		},
		{
			description: "invalid empty value in list",
			value1:      utils.Ptr(fmt.Sprintf(",%s", testUUID1)),
			isValid:     false,
		},
		{
			description:   "no values",
			expectedValue: nil,
			isValid:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := UUIDSliceFlag()
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			var err1, err2 error
			if tt.value1 != nil {
				err1 = cmd.Flags().Set("test-flag", *tt.value1)
			}
			if tt.value2 != nil {
				err2 = cmd.Flags().Set("test-flag", *tt.value2)
			}

			if !tt.isValid && err1 == nil && err2 == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err1 != nil {
				t.Fatalf("failed on valid input: %v", err1)
			}
			if err2 != nil {
				t.Fatalf("failed on valid input: %v", err2)
			}
			value := FlagToStringSliceValue(cmd, "test-flag")
			if !reflect.DeepEqual(tt.expectedValue, value) {
				t.Fatalf("flag did not return set value (expected %s, got %s)", tt.expectedValue, value)
			}
		})
	}
}

func TestCIDRFlag(t *testing.T) {
	tests := []struct {
		description string
		value       string
		isValid     bool
	}{
		{
			description: "valid IPv4 block",
			value:       "198.51.100.14/24",
			isValid:     true,
		},
		{
			description: "valid IPv4 block 2",
			value:       "111.222.111.222/22",
			isValid:     true,
		},
		{
			description: "valid IPv4 single",
			value:       "198.51.100.14/32",
			isValid:     true,
		},
		{
			description: "valid IPv4 entire internet",
			value:       "0.0.0.0/0",
			isValid:     true,
		},
		{
			description: "invalid IPv4 block",
			value:       "198.51.100.14/33",
			isValid:     false,
		},
		{
			description: "invalid IPv4 no block",
			value:       "111.222.111.222",
			isValid:     false,
		},
		{
			description: "valid IPv6 block",
			value:       "2001:db8::/48",
			isValid:     true,
		},
		{
			description: "valid IPv6 single",
			value:       "2001:0db8:85a3:08d3::0370:7344/128",
			isValid:     true,
		},
		{
			description: "valid IPv6 entire internet",
			value:       "::/0",
			isValid:     true,
		},
		{
			description: "invalid IPv6 block",
			value:       "2001:0db8:85a3:08d3::0370:7344/129",
			isValid:     false,
		},
		{
			description: "invalid IPv6 no block",
			value:       "2001:0db8:85a3:08d3::0370:7344",
			isValid:     false,
		},
		{
			description: "invalid",
			value:       "invalid-uuid",
			isValid:     false,
		},
		{
			description: "empty",
			value:       "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := CIDRFlag()
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			err := cmd.Flags().Set("test-flag", tt.value)

			if !tt.isValid && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			value := FlagToStringValue(cmd, "test-flag")
			if value != tt.value {
				t.Fatalf("flag did not return set value")
			}
		})
	}
}

func TestCIDRSliceFlag(t *testing.T) {
	tests := []struct {
		description   string
		value1        *string
		value2        *string
		expectedValue []string
		isValid       bool
	}{
		{
			description:   "valid two single values",
			value1:        utils.Ptr("198.51.100.14/24"),
			value2:        utils.Ptr("198.51.100.14/32"),
			expectedValue: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:       true,
		},
		{
			description:   "valid list value",
			value1:        utils.Ptr("198.51.100.14/24,198.51.100.14/32"),
			expectedValue: []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:       true,
		},
		{
			description:   "valid list value and single value",
			value1:        utils.Ptr("198.51.100.14/24,198.51.100.14/32"),
			value2:        utils.Ptr("111.222.111.222/22"),
			expectedValue: []string{"198.51.100.14/24", "198.51.100.14/32", "111.222.111.222/22"},
			isValid:       true,
		},
		{
			description:   "valid two list values",
			value1:        utils.Ptr("198.51.100.14/24,198.51.100.14/32"),
			value2:        utils.Ptr("111.222.111.222/22,2001:db8::/48"),
			expectedValue: []string{"198.51.100.14/24", "198.51.100.14/32", "111.222.111.222/22", "2001:db8::/48"},
			isValid:       true,
		},
		{
			description: "invalid value",
			value1:      utils.Ptr("invalid-cidr"),
			value2:      utils.Ptr("198.51.100.14/24"),
			isValid:     false,
		},
		{
			description: "invalid value in list",
			value1:      utils.Ptr("198.51.100.14/24,invalid-cidr"),
			isValid:     false,
		},
		{
			description: "invalid empty value",
			value1:      utils.Ptr(""),
			isValid:     false,
		},
		{
			description: "invalid empty value in list",
			value1:      utils.Ptr("198.51.100.14/24,198.51.100.14/32,"),
			isValid:     false,
		},
		{
			description:   "no values",
			expectedValue: nil,
			isValid:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			flag := CIDRSliceFlag()
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			var err1, err2 error
			if tt.value1 != nil {
				err1 = cmd.Flags().Set("test-flag", *tt.value1)
			}
			if tt.value2 != nil {
				err2 = cmd.Flags().Set("test-flag", *tt.value2)
			}

			if !tt.isValid && err1 == nil && err2 == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}

			if err1 != nil {
				t.Fatalf("failed on valid input: %v", err1)
			}
			if err2 != nil {
				t.Fatalf("failed on valid input: %v", err2)
			}
			value := FlagToStringSliceValue(cmd, "test-flag")
			if !reflect.DeepEqual(tt.expectedValue, value) {
				t.Fatalf("flag did not return set value (expected %s, got %s)", tt.expectedValue, value)
			}
		})
	}
}

func TestReadFromFileFlag(t *testing.T) {
	tests := []struct {
		description      string
		value            string
		fileValue        string
		readFileFails    bool
		expectedValue    string
		expectedFileName string // If "", file isn't exected to be called
	}{
		{
			description:   "base",
			value:         "foo",
			expectedValue: "foo",
		},
		{
			description:   "empty",
			value:         "",
			expectedValue: "",
		},
		{
			description:      "read from file 1",
			value:            "@foo",
			fileValue:        "bar",
			expectedValue:    "bar",
			expectedFileName: "foo",
		},
		{
			description:      "read from file 2",
			value:            "@\"foo\"",
			fileValue:        "bar",
			expectedValue:    "bar",
			expectedFileName: "foo",
		},
		{
			description:      "read from file 3",
			value:            "@'foo'",
			fileValue:        "bar",
			expectedValue:    "bar",
			expectedFileName: "foo",
		},
		{
			description:      "read from file fails",
			value:            "@foo",
			readFileFails:    true,
			expectedFileName: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			readFileCalled := false
			reader := func(filename string) ([]byte, error) {
				readFileCalled = true
				if filename != tt.expectedFileName {
					t.Errorf("expected file name %q, got %q instead", tt.expectedFileName, filename)
				}
				if tt.readFileFails {
					return nil, fmt.Errorf("something failed")
				}
				return []byte(tt.fileValue), nil
			}

			flag := &readFromFileFlag{
				reader: reader,
			}
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}
			cmd.Flags().Var(flag, "test-flag", "test")

			err := cmd.Flags().Set("test-flag", tt.value)

			if !readFileCalled && (tt.expectedFileName != "") {
				t.Errorf("read file should've been called")
			}
			if readFileCalled && (tt.expectedFileName == "") {
				t.Errorf("read file shouldn't have been called")
			}
			if tt.readFileFails && err == nil {
				t.Fatalf("did not fail on invalid input")
			}
			if tt.readFileFails {
				return
			}

			if err != nil {
				t.Fatalf("failed on valid input: %v", err)
			}
			value := FlagToStringValue(cmd, "test-flag")
			if value != tt.expectedValue {
				t.Fatalf("flag returned %q, expected %q", value, tt.expectedValue)
			}
		})
	}
}
