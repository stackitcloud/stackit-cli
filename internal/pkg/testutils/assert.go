// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package testutils

// Package test provides utilities for validating CLI command test results with
// explicit helpers for error expectations and value comparisons. By splitting
// error and value handling the package keeps assertions simple and removes the
// need for dynamic type checks in every test case.
//
// Example usage:
//
//	// Expect a specific error type
//	if !test.AssertError(t, run(), &cliErr.FlagValidationError{}) {
//		return
//	}
//
//	// Expect any error
//	if !test.AssertError(t, run(), true) {
//		return
//	}
//
//	// Expect error message substring
//	if !test.AssertError(t, run(), "not found") {
//		return
//	}
//
//	// Compare complex structs with private fields
//	test.AssertValue(t, got, want, test.WithAllowUnexported(MyStruct{}))

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// AssertError verifies that an observed error satisfies the expected condition.
//
// Returns:
//   - bool: True if the test should continue to value checks (i.e., no error occurred).
//
// Behavior:
//  1. If err is nil:
//     - If want is nil or false: Success.
//     - If want is anything else: Fails test (Expected error but got nil).
//  2. If err is non-nil:
//     - If want is nil or false: Fails test (Unexpected error).
//     - If want is true: Success (Any error accepted).
//     - If want is string: Asserts err.Error() contains the string.
//     - If want is error: Asserts errors.Is(err, want) or type match.
func AssertError(t testing.TB, got error, want any) bool {
	t.Helper()

	// Case 1: No error occurred
	if got == nil {
		if want == nil || want == false {
			return true
		}
		t.Errorf("got nil error, want %v", want)
		return false
	}

	// Case 2: Error occurred
	if want == nil || want == false {
		t.Errorf("got unexpected error: %v", got)
		return false
	}

	if want == true {
		return false // Error expected and received, stop test
	}

	// Handle string error type expectation
	if wantStr, ok := want.(string); ok {
		if !strings.Contains(got.Error(), wantStr) {
			t.Errorf("got error %q, want substring %q", got, wantStr)
		}
		return false
	}

	// Handle specific error type expectation
	if wantErr, ok := want.(error); ok {
		if checkErrorMatch(got, wantErr) {
			return false
		}
		t.Errorf("got error %v, want %v", got, wantErr)
		return false
	}

	t.Errorf("invalid want type %T for AssertError", want)
	return false
}

func checkErrorMatch(got, want error) bool {
	if errors.Is(got, want) {
		return true
	}

	// Fallback to type check using errors.As to handle wrapped errors
	if want != nil {
		typ := reflect.TypeOf(want)
		// errors.As requires a pointer to the target type.
		// reflect.New(typ) returns *T where T is the type of want.
		target := reflect.New(typ).Interface()
		if errors.As(got, target) {
			return true
		}
	}

	return false
}

// DiffFunc compares two values and returns a diff string. An empty string means
// equality.
type DiffFunc func(got, want any) string

// ValueComparisonOption configures how HandleValueResult applies cmp options or
// diffing strategies.
type ValueComparisonOption func(*valueComparisonConfig)

type valueComparisonConfig struct {
	diffFunc   DiffFunc
	cmpOptions []cmp.Option
}

func (config *valueComparisonConfig) getDiffFunc() DiffFunc {
	if config.diffFunc != nil {
		return config.diffFunc
	}
	return func(got, want any) string {
		return cmp.Diff(got, want, config.cmpOptions...)
	}
}

// WithCmpOptions accumulates cmp.Options used during value comparison.
func WithAssertionCmpOptions(opts ...cmp.Option) ValueComparisonOption {
	return func(config *valueComparisonConfig) {
		config.cmpOptions = append(config.cmpOptions, opts...)
	}
}

// WithAllowUnexported enables comparison of unexported fields for the provided
// struct types.
func WithAllowUnexported(types ...any) ValueComparisonOption {
	return WithAssertionCmpOptions(cmp.AllowUnexported(types...))
}

// WithDiffFunc sets a custom diffing function. Providing this option overrides
// the default cmp-based diff logic.
func WithDiffFunc(diffFunc DiffFunc) ValueComparisonOption {
	return func(config *valueComparisonConfig) {
		config.diffFunc = diffFunc
	}
}

// WithIgnoreFields ignores the specified fields on the provided type during comparison.
// It uses cmpopts.IgnoreFields to ensure type-safe filtering.
func WithIgnoreFields(typ any, names ...string) ValueComparisonOption {
	return WithAssertionCmpOptions(cmpopts.IgnoreFields(typ, names...))
}

// AssertValue compares two values with cmp.Diff while allowing callers to
// tweak the diff strategy via ValueComparisonOption. A non-empty diff is
// reported as an error containing the diff output.
func AssertValue[T any](t testing.TB, got, want T, opts ...ValueComparisonOption) {
	t.Helper()
	// Configure comparison options
	config := &valueComparisonConfig{}
	for _, opt := range opts {
		opt(config)
	}
	// Perform comparison and report diff
	diff := config.getDiffFunc()(got, want)
	if diff != "" {
		t.Errorf("values do not match: %s", diff)
	}
}
