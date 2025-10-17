package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func TestParseInput[T any](t *testing.T, cmdFactory func(*params.CmdParams) *cobra.Command, parseInputFunc func(*print.Printer, *cobra.Command, []string) (T, error), expectedModel T, argValues []string, flagValues map[string]string, isValid bool) {
	TestParseInputWithAdditionalFlags(t, cmdFactory, parseInputFunc, expectedModel, argValues, flagValues, map[string][]string{}, isValid)
}

func TestParseInputWithAdditionalFlags[T any](t *testing.T, cmdFactory func(*params.CmdParams) *cobra.Command, parseInputFunc func(*print.Printer, *cobra.Command, []string) (T, error), expectedModel T, argValues []string, flagValues map[string]string, additionalFlagValues map[string][]string, isValid bool) {
	p := print.NewPrinter()
	cmd := cmdFactory(&params.CmdParams{Printer: p})
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
