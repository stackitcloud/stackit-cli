// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package utils

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
)

// A fixed time for deterministic month calculations
var fixedNow = time.Date(2025, 10, 31, 12, 0, 0, 0, time.UTC)

// Custom converters for edge case testing
type testNegativeConverter struct{}

func (tnc testNegativeConverter) ToDuration(value uint64, _ time.Time) (time.Duration, error) {
	// Directly return a negative duration to test the secondsFloat < 0 case
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("value exceeds MaxInt64")
	}
	return -time.Duration(int64(value)) * time.Second, nil
}

func TestConvertToSeconds(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		opts    []Option
		want    uint64
		wantErr error
	}{
		// Basic Success Cases
		{
			name:    "seconds",
			timeStr: "30s",
			want:    30,
		},
		{
			name:    "minutes",
			timeStr: "2m",
			want:    120,
		},
		{
			name:    "hours",
			timeStr: "1h",
			want:    3600,
		},
		{
			name:    "days",
			timeStr: "1d",
			want:    86400,
		},
		{
			name:    "months",
			timeStr: "1M",
			opts:    []Option{WithNow(fixedNow)},
			want:    uint64(fixedNow.AddDate(0, 1, 0).Sub(fixedNow).Seconds()),
		},
		// Large values that should work
		{
			name:    "large but valid seconds",
			timeStr: "86400s", // 1 day in seconds
			want:    86400,
		},
		{
			name:    "large but valid hours",
			timeStr: "24h",
			want:    86400,
		},
		// Mixed boundary conditions
		{
			name:    "exactly at min and max",
			timeStr: "100s",
			opts:    []Option{WithMinSeconds(100), WithMaxSeconds(100)},
			want:    100,
		},
		{
			name:    "max seconds is 0 (no limit)",
			timeStr: "999999s",
			opts:    []Option{WithMaxSeconds(0)},
			want:    999999,
		},

		{
			name:    "zero multiplier custom unit",
			timeStr: "5z", // Any value with zero multiplier should return 0
			opts: []Option{WithUnits(map[string]DurationConverter{
				"z": fixedMultiplier{Multiplier: 0},
			})},
			want: 0,
		},
		{
			name:    "zero value",
			timeStr: "0s",
			wantErr: &ValidationError{Type: ValidationErrorInvalidValue},
		},
		{
			name:    "negative multiplier custom unit",
			timeStr: "5n", // Any value with negative multiplier should error
			opts: []Option{WithUnits(map[string]DurationConverter{
				"n": fixedMultiplier{Multiplier: -time.Second},
			})},
			wantErr: &CalculationError{Type: CalculationErrorNegativeMultiplier},
		},
		{
			name:    "result exceeds MaxUint64",
			timeStr: fmt.Sprintf("%dx", math.MaxUint32),
			opts: []Option{WithUnits(map[string]DurationConverter{
				"x": fixedMultiplier{Multiplier: time.Duration(math.MaxInt64)}, // Very large multiplier
			})},
			wantErr: &CalculationError{Type: CalculationErrorOutOfBounds},
		},
		{
			name:    "negative result from calculation",
			timeStr: "5neg", // Use custom converter that returns negative duration
			opts: []Option{WithUnits(map[string]DurationConverter{
				"neg": testNegativeConverter{}, // Custom converter that returns negative duration
			})},
			wantErr: &CalculationError{Type: CalculationErrorNegativeResult},
		},

		// Month edge cases (calendar-aware)
		{
			name:    "month from end of month",
			timeStr: "1M",
			opts:    []Option{WithNow(time.Date(2025, 1, 31, 12, 0, 0, 0, time.UTC))}, // Jan 31 -> Feb 28/29
			want:    uint64(time.Date(2025, 3, 3, 12, 0, 0, 0, time.UTC).Sub(time.Date(2025, 1, 31, 12, 0, 0, 0, time.UTC)).Seconds()),
		},
		{
			name:    "multiple months",
			timeStr: "3M",
			opts:    []Option{WithNow(fixedNow)},
			want:    uint64(fixedNow.AddDate(0, 3, 0).Sub(fixedNow).Seconds()),
		},
		{
			name:    "month value too large for MaxInt",
			timeStr: fmt.Sprintf("%dM", uint64(math.MaxInt)+1),
			opts:    []Option{WithNow(fixedNow)},
			wantErr: &CalculationError{Type: CalculationErrorOutOfBounds},
		},

		// Boundary Checks (min/max)
		{
			name:    "below minimum",
			timeStr: "59s",
			opts:    []Option{WithMinSeconds(60)},
			wantErr: &ValidationError{Type: ValidationErrorBelowMinimum},
		},
		{
			name:    "at minimum",
			timeStr: "60s",
			opts:    []Option{WithMinSeconds(60)},
			want:    60,
		},
		{
			name:    "above maximum",
			timeStr: "61s",
			opts:    []Option{WithMaxSeconds(60)},
			wantErr: &ValidationError{Type: ValidationErrorAboveMaximum},
		},
		{
			name:    "at maximum",
			timeStr: "60s",
			opts:    []Option{WithMaxSeconds(60)},
			want:    60,
		},
		{
			name:    "within boundaries",
			timeStr: "30s",
			opts:    []Option{WithMinSeconds(10), WithMaxSeconds(40)},
			want:    30,
		},

		// Custom units
		{
			name:    "custom unit",
			timeStr: "2w",
			opts: []Option{WithUnits(map[string]DurationConverter{
				"w": fixedMultiplier{Multiplier: 7 * 24 * time.Hour},
			})},
			want: 1209600, // 2 * 7 * 24 * 3600
		},
		{
			name:    "custom multi char unit",
			timeStr: "3wk",
			opts: []Option{WithUnits(map[string]DurationConverter{
				"wk": fixedMultiplier{Multiplier: 7 * 24 * time.Hour},
			})},
			want: 1814400, // 3 * 7 * 24 * 3600
		},
		// Whitespace handling
		{
			name:    "leading whitespace",
			timeStr: " 10s",
			want:    10,
		},
		{
			name:    "trailing whitespace",
			timeStr: "10s ",
			want:    10,
		},

		// Leading zeros (should be invalid)
		{
			name:    "leading zeros invalid",
			timeStr: "01s",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "multiple leading zeros",
			timeStr: "007m",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},

		// Other error cases
		{
			name:    "invalid format no unit",
			timeStr: "123",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "invalid format starts with unit",
			timeStr: "m30",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "invalid format empty string",
			timeStr: "",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "invalid format value missing",
			timeStr: "abcS",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "unsupported unit",
			timeStr: "1y",
			wantErr: &ValidationError{Type: ValidationErrorInvalidUnit},
		},
		{
			name:    "multi-char unit with default units",
			timeStr: "5ms",
			wantErr: &ValidationError{Type: ValidationErrorInvalidUnit},
		},
		{
			name:    "value too large for int64 duration",
			timeStr: fmt.Sprintf("%ds", math.MaxInt64), // This will overflow
			wantErr: &CalculationError{Type: CalculationErrorOutOfBounds},
		},
		{
			name:    "very large number string exceeds uint64 when parsing",
			timeStr: "999999999999999999999s", // Larger than uint64 max
			wantErr: &ValidationError{Type: ValidationErrorInvalidValue},
		},
		// Behavior with non-integer and negative values
		{
			name:    "floating point value",
			timeStr: "1.5h",
			wantErr: &ValidationError{Type: ValidationErrorInvalidValue},
		},
		{
			name:    "negative value",
			timeStr: "-10s",
			wantErr: &ValidationError{Type: ValidationErrorInvalidFormat},
		},
		{
			name:    "comma as decimal separator",
			timeStr: "1,5h",
			wantErr: &ValidationError{Type: ValidationErrorInvalidValue},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var seconds uint64
			var err error

			seconds, err = ConvertToSeconds(tt.timeStr, tt.opts...)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error '%v', but got nil (%v)", tt.wantErr, tt.name)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error to be '%v', but got '%v' (%v)", tt.wantErr, err, tt.name)
				}
				return // Test passed
			}

			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			if seconds != tt.want {
				t.Errorf("expected %d seconds, but got %d", tt.want, seconds)
			}
		})
	}
}

func TestValidationErrorString(t *testing.T) {
	tests := []struct {
		name string
		err  *ValidationError
		want string
	}{
		// InvalidFormat cases
		{
			name: "invalid format with input and reason",
			err:  &ValidationError{Type: ValidationErrorInvalidFormat, Input: "30m", Reason: "leading zeros not allowed"},
			want: `invalid time string format "30m": leading zeros not allowed`,
		},
		{
			name: "invalid format with input only",
			err:  &ValidationError{Type: ValidationErrorInvalidFormat, Input: "30m"},
			want: `invalid time string format: "30m"`,
		},
		{
			name: "invalid format minimal",
			err:  &ValidationError{Type: ValidationErrorInvalidFormat},
			want: "invalid time string format",
		},

		// InvalidValue cases
		{
			name: "invalid value with input and reason",
			err:  &ValidationError{Type: ValidationErrorInvalidValue, Input: "abc", Reason: "not a number"},
			want: `invalid time value "abc": not a number`,
		},
		{
			name: "invalid value with input only",
			err:  &ValidationError{Type: ValidationErrorInvalidValue, Input: "abc"},
			want: `invalid time value: "abc"`,
		},
		{
			name: "invalid value minimal",
			err:  &ValidationError{Type: ValidationErrorInvalidValue},
			want: "invalid time value",
		},

		// InvalidUnit cases
		{
			name: "invalid unit with valid units list",
			err: &ValidationError{
				Type:    ValidationErrorInvalidUnit,
				Input:   "x",
				Context: map[string]any{"validUnits": []string{"s", "m", "h"}},
			},
			want: `invalid time unit "x", supported units are [s m h]`,
		},
		{
			name: "invalid unit with input only",
			err:  &ValidationError{Type: ValidationErrorInvalidUnit, Input: "x"},
			want: `invalid time unit: "x"`,
		},
		{
			name: "invalid unit minimal",
			err:  &ValidationError{Type: ValidationErrorInvalidUnit},
			want: "invalid time unit",
		},

		// BelowMinimum cases
		{
			name: "below minimum with values",
			err: &ValidationError{
				Type:    ValidationErrorBelowMinimum,
				Context: map[string]any{"value": uint64(50), "minimum": uint64(60)},
			},
			want: "duration is below minimum: 50 seconds (minimum: 60 seconds)",
		},
		{
			name: "below minimum minimal",
			err:  &ValidationError{Type: ValidationErrorBelowMinimum},
			want: "duration is below the allowed minimum",
		},

		// AboveMaximum cases
		{
			name: "above maximum with values",
			err: &ValidationError{
				Type:    ValidationErrorAboveMaximum,
				Context: map[string]any{"value": uint64(120), "maximum": uint64(100)},
			},
			want: "duration exceeds maximum: 120 seconds (maximum: 100 seconds)",
		},
		{
			name: "above maximum minimal",
			err:  &ValidationError{Type: ValidationErrorAboveMaximum},
			want: "duration exceeds the allowed maximum",
		},

		// Default case
		{
			name: "unknown validation error type",
			err:  &ValidationError{Type: "unknown_type", Input: "test"},
			want: "validation error: test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("ValidationError.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalculationErrorString(t *testing.T) {
	tests := []struct {
		name string
		err  *CalculationError
		want string
	}{
		// OutOfBounds cases
		{
			name: "out of bounds with all context",
			err: &CalculationError{
				Type:    CalculationErrorOutOfBounds,
				Value:   12345,
				Context: map[string]any{"operation": "multiplication", "limit": "MaxInt64"},
			},
			want: "calculation result is out of bounds (value: 12345) during multiplication (exceeds MaxInt64)",
		},
		{
			name: "out of bounds with value and operation",
			err: &CalculationError{
				Type:    CalculationErrorOutOfBounds,
				Value:   12345,
				Context: map[string]any{"operation": "multiplication"},
			},
			want: "calculation result is out of bounds (value: 12345) during multiplication",
		},
		{
			name: "out of bounds with value only",
			err: &CalculationError{
				Type:  CalculationErrorOutOfBounds,
				Value: 12345,
			},
			want: "calculation result is out of bounds (value: 12345)",
		},
		{
			name: "out of bounds minimal",
			err:  &CalculationError{Type: CalculationErrorOutOfBounds},
			want: "calculation result is out of bounds",
		},

		// NegativeResult cases
		{
			name: "negative result with result value",
			err: &CalculationError{
				Type:    CalculationErrorNegativeResult,
				Context: map[string]any{"result": -123.456},
			},
			want: "calculated duration is negative: -123.456000",
		},
		{
			name: "negative result minimal",
			err:  &CalculationError{Type: CalculationErrorNegativeResult},
			want: "calculated duration is negative",
		},

		// NegativeMultiplier cases
		{
			name: "negative multiplier with multiplier value",
			err: &CalculationError{
				Type:    CalculationErrorNegativeMultiplier,
				Context: map[string]any{"multiplier": -5 * time.Second},
			},
			want: "duration multiplier is negative: -5s",
		},
		{
			name: "negative multiplier minimal",
			err:  &CalculationError{Type: CalculationErrorNegativeMultiplier},
			want: "duration multiplier is negative",
		},

		// Default case
		{
			name: "unknown calculation error type",
			err:  &CalculationError{Type: "unknown_type", Value: 123},
			want: "calculation error with value 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("CalculationError.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidationErrorIs(t *testing.T) {
	tests := []struct {
		name   string
		err    *ValidationError
		target error
		want   bool
	}{
		// True cases - same type and same error type
		{
			name:   "same validation error type matches",
			err:    &ValidationError{Type: ValidationErrorInvalidFormat},
			target: &ValidationError{Type: ValidationErrorInvalidFormat},
			want:   true,
		},
		{
			name:   "different validation error types don't match",
			err:    &ValidationError{Type: ValidationErrorInvalidFormat},
			target: &ValidationError{Type: ValidationErrorInvalidValue},
			want:   false,
		},
		// False cases - different error types
		{
			name:   "validation error vs calculation error",
			err:    &ValidationError{Type: ValidationErrorInvalidFormat},
			target: &CalculationError{Type: CalculationErrorOutOfBounds},
			want:   false,
		},
		{
			name:   "validation error vs standard error",
			err:    &ValidationError{Type: ValidationErrorInvalidFormat},
			target: errors.New("some other error"),
			want:   false,
		},
		{
			name:   "validation error vs nil",
			err:    &ValidationError{Type: ValidationErrorInvalidFormat},
			target: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Is(tt.target)
			if got != tt.want {
				t.Errorf("ValidationError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculationErrorIs(t *testing.T) {
	tests := []struct {
		name   string
		err    *CalculationError
		target error
		want   bool
	}{
		// True cases - same type and same error type
		{
			name:   "same calculation error type matches",
			err:    &CalculationError{Type: CalculationErrorOutOfBounds},
			target: &CalculationError{Type: CalculationErrorOutOfBounds},
			want:   true,
		},
		// False cases - different calculation error types
		{
			name:   "different calculation error types don't match",
			err:    &CalculationError{Type: CalculationErrorOutOfBounds},
			target: &CalculationError{Type: CalculationErrorNegativeResult},
			want:   false,
		},
		// False cases - different error types
		{
			name:   "calculation error vs validation error",
			err:    &CalculationError{Type: CalculationErrorOutOfBounds},
			target: &ValidationError{Type: ValidationErrorInvalidFormat},
			want:   false,
		},
		{
			name:   "calculation error vs standard error",
			err:    &CalculationError{Type: CalculationErrorOutOfBounds},
			target: errors.New("some other error"),
			want:   false,
		},
		{
			name:   "calculation error vs nil",
			err:    &CalculationError{Type: CalculationErrorOutOfBounds},
			target: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Is(tt.target)
			if got != tt.want {
				t.Errorf("CalculationError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
