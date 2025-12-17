package template

import (
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
)

var (
	testProjectId = uuid.NewString()
	testRegion    = "eu01"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectId, Region: testRegion, Verbosity: globalflags.VerbosityDefault},
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "alb with yaml",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[formatFlag] = "yaml"
				flagValues[typeFlag] = "alb"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Format = utils.Ptr("yaml")
				model.Type = utils.Ptr("alb")
			}),
		}, {
			description: "alb with yaml",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[formatFlag] = "yaml"
				flagValues[typeFlag] = "alb"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Format = utils.Ptr("yaml")
				model.Type = utils.Ptr("alb")
			}),
		}, {
			description: "alb with json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[formatFlag] = "json"
				flagValues[typeFlag] = "alb"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Format = utils.Ptr("json")
				model.Type = utils.Ptr("alb")
			}),
		}, {
			description: "pool with yaml",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[formatFlag] = "yaml"
				flagValues[typeFlag] = "pool"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Format = utils.Ptr("yaml")
				model.Type = utils.Ptr("pool")
			}),
		},
		{
			description: "pool with json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[formatFlag] = "json"
				flagValues[typeFlag] = "pool"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Format = utils.Ptr("json")
				model.Type = utils.Ptr("pool")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}
