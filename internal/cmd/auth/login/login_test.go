package login

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		portFlag:          "8010",
		useDeviceFlowFlag: "false",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		Port:          utils.Ptr(8010),
		UseDeviceFlow: false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		argValues     []string
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: &inputModel{
				Port:          nil,
				UseDeviceFlow: false,
			},
		},
		{
			description: "device flow enabled",
			flagValues: map[string]string{
				useDeviceFlowFlag: "true",
			},
			isValid: true,
			expectedModel: &inputModel{
				Port:          nil,
				UseDeviceFlow: true,
			},
		},
		{
			description: "lower limit",
			flagValues: map[string]string{
				portFlag: "8000",
			},
			isValid: true,
			expectedModel: &inputModel{
				Port:          utils.Ptr(8000),
				UseDeviceFlow: false,
			},
		},
		{
			description: "below lower limit is not valid ",
			flagValues: map[string]string{
				portFlag: "7999",
			},
			isValid: false,
		},
		{
			description: "upper limit",
			flagValues: map[string]string{
				portFlag: "8020",
			},
			isValid: true,
			expectedModel: &inputModel{
				Port:          utils.Ptr(8020),
				UseDeviceFlow: false,
			},
		},
		{
			description: "above upper limit is not valid ",
			flagValues: map[string]string{
				portFlag: "8021",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}
