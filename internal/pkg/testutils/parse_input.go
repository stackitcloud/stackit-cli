// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package testutils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

// ParseInputTestCase aggregates all required elements to exercise a CLI parseInput
// function. It centralizes the common flag setup, validation, and result
// assertions used throughout the edge command test suites.
type ParseInputTestCase[T any] struct {
	Name string
	// Args simulates positional arguments passed to the command.
	Args []string
	// Flags sets simple single-value flags.
	Flags map[string]string
	// RepeatFlags sets flags that can be specified multiple times (e.g. slice flags).
	RepeatFlags map[string][]string
	WantModel   T
	WantErr     any
	CmdFactory  func(*types.CmdParams) *cobra.Command
	// ParseInputFunc is the function under test. It must accept the printer, command, and args.
	ParseInputFunc func(*print.Printer, *cobra.Command, []string) (T, error)
}

// ParseInputCaseOption allows configuring the test execution behavior.
type ParseInputCaseOption func(*parseInputCaseConfig)

type parseInputCaseConfig struct {
	cmpOpts []ValueComparisonOption
}

// WithParseInputCmpOptions sets custom comparison options for AssertValue.
func WithParseInputCmpOptions(opts ...ValueComparisonOption) ParseInputCaseOption {
	return func(cfg *parseInputCaseConfig) {
		cfg.cmpOpts = append(cfg.cmpOpts, opts...)
	}
}

func defaultParseInputCaseConfig() *parseInputCaseConfig {
	return &parseInputCaseConfig{}
}

// RunParseInputCase executes a single parse-input test case using the provided
// configuration. It mirrors the typical table-driven pattern while removing the
// boilerplate repeated across tests. The helper short-circuits as soon as an
// expected error is encountered.
func RunParseInputCase[T any](t *testing.T, tc ParseInputTestCase[T], opts ...ParseInputCaseOption) {
	t.Helper()

	cfg := defaultParseInputCaseConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if tc.CmdFactory == nil {
		t.Fatalf("parse input case %q missing CmdFactory", tc.Name)
	}
	if tc.ParseInputFunc == nil {
		t.Fatalf("parse input case %q missing ParseInputFunc", tc.Name)
	}

	printer := print.NewPrinter()
	cmd := tc.CmdFactory(&types.CmdParams{Printer: printer})
	if cmd == nil {
		t.Fatalf("parse input case %q produced nil command", tc.Name)
	}
	if printer.Cmd == nil {
		printer.Cmd = cmd
	}

	if err := globalflags.Configure(cmd.Flags()); err != nil {
		t.Fatalf("configure global flags: %v", err)
	}

	// Set regular flag values.
	for flag, value := range tc.Flags {
		if err := cmd.Flags().Set(flag, value); err != nil {
			AssertError(t, err, tc.WantErr)
			return
		}
	}

	// Set repeated flag values.
	for flag, values := range tc.RepeatFlags {
		for _, value := range values {
			if err := cmd.Flags().Set(flag, value); err != nil {
				AssertError(t, err, tc.WantErr)
				return
			}
		}
	}

	// Test cobra argument validation.
	if err := cmd.ValidateArgs(tc.Args); err != nil {
		AssertError(t, err, tc.WantErr)
		return
	}

	// Test cobra required flags validation.
	if err := cmd.ValidateRequiredFlags(); err != nil {
		AssertError(t, err, tc.WantErr)
		return
	}

	// Test cobra flag group validation.
	if err := cmd.ValidateFlagGroups(); err != nil {
		AssertError(t, err, tc.WantErr)
		return
	}

	// Test parse input function.
	got, err := tc.ParseInputFunc(printer, cmd, tc.Args)
	if !AssertError(t, err, tc.WantErr) {
		return
	}

	AssertValue(t, got, tc.WantModel, cfg.cmpOpts...)
}
