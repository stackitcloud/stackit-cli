package export

import (
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/google/go-cmp/cmp"
)

const (
	testProfileArg = "default"
	testExportPath = "/tmp/stackit-profiles/" + testProfileArg + ".json"
)

func fixtureArgValues(mods ...func(args []string)) []string {
	args := []string{
		testProfileArg,
	}
	for _, mod := range mods {
		mod(args)
	}
	return args
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		filePathFlag: testExportPath,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(inputModel *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		ProfileName: testProfileArg,
		FilePath:    testExportPath,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argsValues    []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argsValues:    fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argsValues:  []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no args",
			argsValues:  []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flags",
			argsValues:  fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: fixtureInputModel(func(inputModel *inputModel) {
				inputModel.FilePath = fmt.Sprintf("%s.json", testProfileArg)
			}),
		},
		{
			description: "custom file-path without file extension",
			argsValues:  fixtureArgValues(),
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					flagValues[filePathFlag] = "./my-exported-config"
				}),
			isValid: true,
			expectedModel: fixtureInputModel(func(inputModel *inputModel) {
				inputModel.FilePath = "./my-exported-config"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err = cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateArgs(tt.argsValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argsValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
