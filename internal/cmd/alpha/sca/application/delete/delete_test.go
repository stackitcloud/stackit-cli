package delete

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	testRegion = "eu01"
)

var testClient = &sca.APIClient{DefaultAPI: &sca.DefaultAPIService{}}
var testCtx = context.Background()

var (
	testProjectID     = uuid.NewString()
	testEnvironmentID = uuid.NewString()
	testApplicationID = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testApplicationID,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,

		environmentIDFlag: testEnvironmentID,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectID,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		EnvironmentID: testEnvironmentID,
		ApplicationID: testApplicationID,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sca.ApiDeleteApplicationRequest)) sca.ApiDeleteApplicationRequest {
	request := testClient.DefaultAPI.DeleteApplication(testCtx, testProjectID, testEnvironmentID, testApplicationID)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		desc          string
		argValues     []string
		flagValues    map[string]string
		expectedModel *inputModel
		isValid       bool
	}{
		{
			desc:          "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			expectedModel: fixtureInputModel(),
			isValid:       true,
		},
		{
			desc:      "required only",
			argValues: fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, environmentIDFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.EnvironmentID = testProjectID
			}),
			isValid: true,
		},
		{
			desc:       "no arg values",
			argValues:  []string{},
			flagValues: fixtureFlagValues(),
			isValid:    false,
		},
		{
			desc:       "no values",
			argValues:  nil,
			flagValues: nil,
			isValid:    false,
		},
		{
			desc:      "project id missing",
			argValues: fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			desc:      "project id invalid",
			argValues: fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc:      "environment id invalid",
			argValues: fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[environmentIDFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sca.ApiDeleteApplicationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, sca.DefaultAPIService{}),
				cmpopts.SortSlices(func(a, b sca.EnvVar) bool {
					return a.Key < b.Key
				}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
