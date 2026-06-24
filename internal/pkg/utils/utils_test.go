package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func TestConvertInt64PToFloat64P(t *testing.T) {
	tests := []struct {
		name     string
		input    *int64
		expected *float64
	}{
		{
			name:     "positive",
			input:    Int64Ptr(1),
			expected: Float64Ptr(1.0),
		},
		{
			name:     "negative",
			input:    Int64Ptr(-1),
			expected: Float64Ptr(-1.0),
		},
		{
			name:     "zero",
			input:    Int64Ptr(0),
			expected: Float64Ptr(0.0),
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := ConvertInt64PToFloat64P(tt.input)

			if expected == nil && tt.expected == nil && tt.input == nil {
				return
			}

			if *expected != *tt.expected {
				t.Errorf("ConvertInt64ToFloat64() = %v, want %v", *expected, *tt.expected)
			}
		})
	}
}

func TestConvertInt32PToFloat64P(t *testing.T) {
	tests := []struct {
		name     string
		input    *int32
		expected *float32
	}{
		{
			name:     "positive",
			input:    utils.Ptr(int32(1)),
			expected: utils.Ptr(float32(1)),
		},
		{
			name:     "negative",
			input:    utils.Ptr(int32(-1)),
			expected: utils.Ptr(float32(-1)),
		},
		{
			name:     "zero",
			input:    utils.Ptr(int32(0)),
			expected: utils.Ptr(float32(0)),
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := ConvertInt32PToFloat32P(tt.input)

			if expected == nil && tt.expected == nil && tt.input == nil {
				return
			}

			if *expected != *tt.expected {
				t.Errorf("ConvertInt64ToFloat64() = %v, want %v", *expected, *tt.expected)
			}
		})
	}
}

func TestConvertInt32PToFloat64PLossyConversion(t *testing.T) {
	i := int32(900_000_001)
	f := ConvertInt32PToFloat32P(&i)
	s := fmt.Sprintf("%f", *f)
	if s != "900000000.000000" {
		t.Errorf("Expected lossy conversion of %d to %f, got %s", i, *f, s)
	}
}

func TestValidateURLDomain(t *testing.T) {
	tests := []struct {
		name             string
		allowedUrlDomain string
		isValid          bool
		input            string
	}{
		{
			name:             "STACKIT URL valid",
			allowedUrlDomain: "stackit.cloud",
			input:            "https://example.stackit.cloud",
			isValid:          true,
		},
		{
			name:             "STACKIT URL invalid",
			allowedUrlDomain: "example.com",
			input:            "https://example.stackit.cloud",
			isValid:          false,
		},
		{
			name:             "non-STACKIT URL invalid",
			allowedUrlDomain: "stackit.cloud",
			input:            "https://www.very-suspicious-website.com/",
			isValid:          false,
		},
		{
			name:             "non-STACKIT URL valid",
			allowedUrlDomain: "example.com",
			input:            "https://www.test.example.com/",
			isValid:          true,
		},
		{
			name:             "every URL valid",
			allowedUrlDomain: "",
			input:            "https://www.test.example.com/",
			isValid:          true,
		},
		{
			name:    "invalid URL",
			input:   "",
			isValid: false,
		},
		{
			name:    "invalid protocol",
			input:   "http://example.stackit.cloud",
			isValid: false,
		},
		{
			name:    "no protocol",
			input:   "example.stackit.cloud",
			isValid: false,
		},
		{
			name:    "valid endpoint",
			input:   "https://service-account.api.stackit.cloud/token",
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			viper.Set(config.AllowedUrlDomainKey, tt.allowedUrlDomain)

			err := ValidateURLDomain(tt.input)
			if tt.isValid && err != nil {
				t.Errorf("expected URL to be valid, got error: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("expected URL to be invalid, got no error")
			}
		})
	}
}

func TestUserAgentConfigOption(t *testing.T) {
	type args struct {
		providerVersion string
	}
	tests := []struct {
		name string
		args args
		want sdkConfig.ConfigurationOption
	}{
		{
			name: "TestUserAgentConfigOption",
			args: args{
				providerVersion: "1.0.0",
			},
			want: sdkConfig.WithUserAgent("stackit-cli/1.0.0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfigActual := sdkConfig.Configuration{}
			err := tt.want(&clientConfigActual)
			if err != nil {
				t.Errorf("error configuring client: %v", err)
			}

			clientConfigExpected := sdkConfig.Configuration{}
			err = UserAgentConfigOption(tt.args.providerVersion)(&clientConfigExpected)
			if err != nil {
				t.Errorf("error configuring client: %v", err)
			}

			if !reflect.DeepEqual(clientConfigActual, clientConfigExpected) {
				t.Errorf("UserAgentConfigOption() = %v, want %v", clientConfigActual, clientConfigExpected)
			}
		})
	}
}

func TestConvertStringMapToInterfaceMap(t *testing.T) {
	tests := []struct {
		name     string
		input    *map[string]string
		expected *map[string]interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty map",
			input:    &map[string]string{},
			expected: nil,
		},
		{
			name: "single key-value pair",
			input: &map[string]string{
				"key1": "value1",
			},
			expected: &map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			name: "multiple key-value pairs",
			input: &map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			expected: &map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "special characters in values",
			input: &map[string]string{
				"key1": "value with spaces",
				"key2": "value,with,commas",
				"key3": "value\nwith\nnewlines",
			},
			expected: &map[string]interface{}{
				"key1": "value with spaces",
				"key2": "value,with,commas",
				"key3": "value\nwith\nnewlines",
			},
		},
		{
			name: "empty values",
			input: &map[string]string{
				"key1": "",
				"key2": "value2",
			},
			expected: &map[string]interface{}{
				"key1": "",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStringMapToInterfaceMap(tt.input)

			// Check if both are nil
			if result == nil && tt.expected == nil {
				return
			}

			// Check if one is nil and other isn't
			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("ConvertStringMapToInterfaceMap() = %v, want %v", result, tt.expected)
				return
			}

			// Compare maps
			if len(*result) != len(*tt.expected) {
				t.Errorf("ConvertStringMapToInterfaceMap() map length = %d, want %d", len(*result), len(*tt.expected))
				return
			}

			for k, v := range *result {
				expectedVal, ok := (*tt.expected)[k]
				if !ok {
					t.Errorf("ConvertStringMapToInterfaceMap() unexpected key %s in result", k)
					continue
				}
				if v != expectedVal {
					t.Errorf("ConvertStringMapToInterfaceMap() value for key %s = %v, want %v", k, v, expectedVal)
				}
			}
		})
	}
}

func TestGetSliceFromPointer(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]string
		expected []string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: []string{},
		},
		{
			name: "pointer to nil slice",
			input: func() *[]string {
				var s []string
				return &s
			}(),
			expected: []string{},
		},
		{
			name:     "empty slice",
			input:    &[]string{},
			expected: []string{},
		},
		{
			name:     "populated slice",
			input:    &[]string{"item1", "item2"},
			expected: []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSliceFromPointer(tt.input)

			if result == nil {
				t.Errorf("GetSliceFromPointer() = %v, want %v", result, tt.expected)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetSliceFromPointer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args[T any, U any] struct {
		input []T
		mapFn func(T) U
	}
	type testCase[T any, U any] struct {
		name string
		args args[T, U]
		want []U
	}
	tests := []testCase[string, *string]{
		{
			name: "default",
			args: args[string, *string]{
				input: []string{"foo", "bar"},
				mapFn: Ptr[string],
			},
			want: []*string{Ptr("foo"), Ptr("bar")},
		},
		{
			name: "input slice is nil",
			args: args[string, *string]{
				input: nil,
				mapFn: Ptr[string],
			},
			want: []*string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.input, tt.args.mapFn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
