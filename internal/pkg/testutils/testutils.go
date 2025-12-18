package testutils

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

// TestParseInput centralizes the logic to test a combination of inputs (arguments, flags) for a cobra command
func TestParseInput[T any](t *testing.T, cmdFactory func(*types.CmdParams) *cobra.Command, parseInputFunc func(*print.Printer, *cobra.Command, []string) (T, error), expectedModel T, argValues []string, flagValues map[string]string, isValid bool) {
	t.Helper()
	TestParseInputWithAdditionalFlags(t, cmdFactory, parseInputFunc, expectedModel, argValues, flagValues, map[string][]string{}, isValid)
}

// TestParseInputWithAdditionalFlags centralizes the logic to test a combination of inputs (arguments, flags) for a cobra command.
// It allows to pass multiple instances of a single flag to the cobra command using the `additionalFlagValues` parameter.
func TestParseInputWithAdditionalFlags[T any](t *testing.T, cmdFactory func(*types.CmdParams) *cobra.Command, parseInputFunc func(*print.Printer, *cobra.Command, []string) (T, error), expectedModel T, argValues []string, flagValues map[string]string, additionalFlagValues map[string][]string, isValid bool) {
	t.Helper()
	p := print.NewPrinter()
	cmd := cmdFactory(&types.CmdParams{Printer: p})
	err := globalflags.Configure(cmd.Flags())
	if err != nil {
		t.Fatalf("configure global flags: %v", err)
	}

	// set regular flag values
	for flag, value := range flagValues {
		err := cmd.Flags().Set(flag, value)
		if err != nil {
			if !isValid {
				return
			}
			t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
		}
	}

	// set additional flag values
	for flag, values := range additionalFlagValues {
		for _, value := range values {
			err := cmd.Flags().Set(flag, value)
			if err != nil {
				if !isValid {
					return
				}
				t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
			}
		}
	}

	if cmd.PreRun != nil {
		// can be used for dynamic flag configuration
		cmd.PreRun(cmd, argValues)
	}

	if cmd.PreRunE != nil {
		err := cmd.PreRunE(cmd, argValues)
		if err != nil {
			if !isValid {
				return
			}
			t.Fatalf("error in PreRunE: %v", err)
		}
	}

	err = cmd.ValidateArgs(argValues)
	if err != nil {
		if !isValid {
			return
		}
		t.Fatalf("error validating args: %v", err)
	}

	err = cmd.ValidateRequiredFlags()
	if err != nil {
		if !isValid {
			return
		}
		t.Fatalf("error validating flags: %v", err)
	}

	err = cmd.ValidateFlagGroups()
	if err != nil {
		if !isValid {
			return
		}
		t.Fatalf("error validating flags: %v", err)
	}

	model, err := parseInputFunc(p, cmd, argValues)
	if err != nil {
		if !isValid {
			return
		}
		t.Fatalf("error parsing input: %v", err)
	}

	if !isValid {
		t.Fatalf("did not fail on invalid input")
	}
	diff := cmp.Diff(model, expectedModel)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}
