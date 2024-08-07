package utils

import (
	"testing"
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

func TestValidateSTACKITURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isValid bool
	}{
		{
			name:    "STACKIT URL",
			input:   "https://example.stackit.cloud",
			isValid: true,
		},
		{
			name:    "non-STACKIT URL",
			input:   "https://www.very-suspicious-website.com/",
			isValid: false,
		},
		{
			name:    "invalid URL",
			input:   "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSTACKITURL(tt.input)
			if tt.isValid && err != nil {
				t.Errorf("expected URL to be valid, got error: %v", err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("expected URL to be invalid, got no error")
			}
		})
	}
}
