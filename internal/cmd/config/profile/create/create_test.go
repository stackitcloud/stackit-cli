package create

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const testProfile = "test-profile"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testProfile,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		Profile:          testProfile,
		FromEmptyProfile: false,
		NoSet:            false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			isValid:     false,
		},
		{
			description: "some global flag",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.VerbosityFlag: globalflags.DebugVerbosity,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Verbosity = globalflags.DebugVerbosity
			}),
		},
		{
			description: "invalid profile",
			argValues:   []string{"invalid-profile-&"},
			isValid:     false,
		},
		{
			description: "use default given",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				fromEmptyProfile: "true",
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.FromEmptyProfile = true
			}),
		},
		{
			description: "no set given",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				noSetFlag: "true",
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NoSet = true
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}
