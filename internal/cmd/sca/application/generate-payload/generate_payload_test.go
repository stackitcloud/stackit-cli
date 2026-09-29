package generatepayload

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
	testFilePath      = "example-file"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,

		environmentIDFlag: testEnvironmentID,
		applicationIDFlag: testApplicationID,
		filePathFlag:      testFilePath,
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
		ApplicationID: &testApplicationID,
		EnvironmentID: &testEnvironmentID,
		FilePath:      &testFilePath,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sca.ApiGetApplicationRequest)) sca.ApiGetApplicationRequest {
	request := testClient.DefaultAPI.GetApplication(testCtx, testProjectID, testEnvironmentID, testApplicationID)
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
			flagValues:    fixtureFlagValues(),
			expectedModel: fixtureInputModel(),
			isValid:       true,
		},
		{
			desc:       "no values",
			flagValues: map[string]string{},
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ProjectId = ""
				model.Region = ""
				model.Verbosity = globalflags.VerbosityDefault
				model.EnvironmentID = nil
				model.ApplicationID = nil
				model.FilePath = nil
			}),
			isValid: true,
		},
		{
			desc: "application id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, applicationIDFlag)
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ApplicationID = nil
			}),
			isValid: true,
		},
		{
			desc: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			desc: "project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc: "environment id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[environmentIDFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			desc: "application id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[applicationIDFlag] = "invalid-uuid"
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
		expectedRequest sca.ApiGetApplicationRequest
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
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
