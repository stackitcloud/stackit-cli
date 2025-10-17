package importProfile

import (
	_ "embed"
	"strconv"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const testProfile = "test-profile"
const testConfig = "@./template/profile.json"
const testNoSet = false

//go:embed template/profile.json
var testConfigContent string

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		nameFlag:   testProfile,
		configFlag: testConfig,
		noSetFlag:  strconv.FormatBool(testNoSet),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
		},
		ProfileName: testProfile,
		Config:      testConfigContent,
		NoSet:       testNoSet,
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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no flags",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "invalid path",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[configFlag] = "@./template/invalid-file"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}
