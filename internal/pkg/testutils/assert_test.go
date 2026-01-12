// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package testutils

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }

type anotherError struct{ code int }

func (e *anotherError) Error() string { return fmt.Sprintf("code=%d", e.code) }

type mockTB struct {
	testing.TB
	failed bool
	msg    string
}

func (m *mockTB) Helper() {}
func (m *mockTB) Errorf(format string, args ...any) {
	m.failed = true
	m.msg = fmt.Sprintf(format, args...)
}

func TestAssertError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("sentinel")

	tests := map[string]struct {
		got     error // The input provided as got to AssertError()
		want    any   // The input provided as want to AssertError()
		wantErr bool  // Whether this comparison is expected to fail
	}{
		"exact match": {
			got:     &customError{msg: "boom"},
			want:    &customError{},
			wantErr: false,
		},
		"error string message match": {
			got:     errors.New("same message"),
			want:    "same message",
			wantErr: false,
		},
		"error string mismatch": {
			got:     errors.New("different"),
			want:    "same message",
			wantErr: true,
		},
		"sentinel via errors.Is": {
			got:     fmt.Errorf("wrap: %w", sentinel),
			want:    sentinel,
			wantErr: false,
		},
		"any error (true)": {
			got:     errors.New("any"),
			want:    true,
			wantErr: false,
		},
		"nil expectation (nil)": {
			got:     nil,
			want:    nil,
			wantErr: false,
		},
		"nil expectation (false)": {
			got:     nil,
			want:    false,
			wantErr: false,
		},
		"nil error input with error expectation": {
			got:     nil,
			want:    true,
			wantErr: true,
		},
		"unexpected error (nil want)": {
			got:     errors.New("unexpected"),
			want:    nil,
			wantErr: true,
		},
		"type match without message": {
			got:     &customError{msg: "alpha"},
			want:    &customError{msg: "beta"},
			wantErr: false,
		},
		"type mismatch": {
			got:     &customError{msg: "alpha"},
			want:    &anotherError{},
			wantErr: true,
		},
		"no error when none expected": {
			got:     nil,
			want:    false,
			wantErr: false,
		},
		"error but want false": {
			got:     errors.New("boom"),
			want:    false,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mock := &mockTB{}
			result := AssertError(mock, tt.got, tt.want)

			// if the test failed but we didn't expect it to fail
			if mock.failed != tt.wantErr {
				t.Fatalf("AssertError() failed = %v, wantErr %v (msg: %s)", mock.failed, tt.wantErr, mock.msg)
			}
			// if we expected an error the result of AssertError() should be false (this is what AssertError() does in case of error)
			if tt.wantErr && result != false {
				t.Fatalf("AssertError() returned = %v, want %v", result, tt.wantErr)
			}
		})
	}
}

func TestCheckErrorMatch(t *testing.T) {
	t.Parallel()

	underlying := &customError{msg: "root"}
	wrapped := fmt.Errorf("wrap: %w", underlying)
	if !checkErrorMatch(wrapped, &customError{}) {
		t.Fatalf("expected wrapped customError to match via errors.As")
	}

	notMatch := errors.New("other")
	if checkErrorMatch(notMatch, &anotherError{}) {
		t.Fatalf("expected mismatch for unrelated error types")
	}
}

func TestAssertValue(t *testing.T) {
	t.Parallel()

	type payload struct {
		Visible string
		hidden  int
	}

	customDiff := func(got, want any) string {
		if reflect.DeepEqual(got, want) {
			return ""
		}
		return "custom diff"
	}

	tests := []struct {
		name    string
		got     any  // The input provided as got to AssertValue()
		want    any  // The input provided as want to AssertValue()
		wantErr bool // Whether this comparison is expected to fail
		opts    []ValueComparisonOption
	}{
		{
			name: "allow unexported success",
			got:  payload{Visible: "ok", hidden: 1},
			want: payload{Visible: "ok", hidden: 1},
			opts: []ValueComparisonOption{WithAllowUnexported(payload{})},
		},
		{
			name:    "allow unexported mismatch",
			got:     payload{Visible: "oops", hidden: 1},
			want:    payload{Visible: "ok", hidden: 1},
			opts:    []ValueComparisonOption{WithAllowUnexported(payload{})},
			wantErr: true,
		},
		{
			name: "cmp options sort",
			got:  []string{"b", "a", "c"},
			want: []string{"a", "b", "c"},
			opts: []ValueComparisonOption{WithAssertionCmpOptions(cmpopts.SortSlices(func(a, b string) bool { return a < b }))},
		},
		{
			name:    "custom diff mismatch",
			got:     1,
			want:    2,
			opts:    []ValueComparisonOption{WithDiffFunc(customDiff)},
			wantErr: true,
		},
		{
			name: "default diff success",
			got:  42,
			want: 42,
		},
		{
			name:    "default diff mismatch",
			got:     1,
			want:    2,
			wantErr: true,
		},
		{
			name: "diff func overrides cmp options",
			got:  []string{"b"},
			want: []string{"a"},
			opts: []ValueComparisonOption{
				WithAssertionCmpOptions(cmpopts.SortSlices(func(a, b string) bool { return a < b })),
				WithDiffFunc(func(_, _ any) string { return "" }),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock := &mockTB{}
			AssertValue(mock, tt.got, tt.want, tt.opts...)

			// if the test failed but we didn't expect it to fail
			if mock.failed != tt.wantErr {
				t.Fatalf("AssertValue failed = %v, want %v (msg: %s)", mock.failed, tt.wantErr, mock.msg)
			}
		})
	}
}
