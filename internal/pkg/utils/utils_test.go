package utils

import (
	"reflect"
	"testing"

	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"

	"github.com/spf13/viper"
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
