package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testContextKey struct{}

const (
	testRegion = "eu01"
	testName   = "example"
)

var (
	testContext    = context.WithValue(context.Background(), testContextKey{}, "token-create")
	testClient     = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectID  = uuid.NewString()
	testInstanceID = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectID,
		globalflags.RegionFlag:    testRegion,
		instanceIDFlag:            testInstanceID,
		nameFlag:                  testName,
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
		InstanceID: testInstanceID,
		Name:       testName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *modelexperiments.ApiCreateInstanceTokenRequest)) modelexperiments.ApiCreateInstanceTokenRequest {
	request := testClient.DefaultAPI.CreateInstanceToken(testContext, testProjectID, testRegion, testInstanceID).
		CreateInstanceTokenPayload(modelexperiments.CreateInstanceTokenPayload{Name: testName})
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "optional fields",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[descriptionFlag] = "service token"
				flagValues[labelsFlag] = "env=prod,team=ml"
				flagValues[ttlDurationFlag] = "5h"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("service token")
				model.Labels = &map[string]string{"env": "prod", "team": "ml"}
				model.TTLDuration = utils.Ptr("5h")
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIDFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = ""
			}),
			isValid: false,
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
			description: "instance id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIDFlag] = "invalid-uuid"
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest modelexperiments.ApiCreateInstanceTokenRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "optional fields",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("service token")
				model.Labels = &map[string]string{"env": "prod", "team": "ml"}
				model.TTLDuration = utils.Ptr("5h")
			}),
			expectedRequest: fixtureRequest(func(request *modelexperiments.ApiCreateInstanceTokenRequest) {
				payload := modelexperiments.CreateInstanceTokenPayload{
					Name:        testName,
					Description: utils.Ptr("service token"),
					Labels:      &map[string]string{"env": "prod", "team": "ml"},
					TtlDuration: utils.Ptr("5h"),
				}
				*request = request.CreateInstanceTokenPayload(payload)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testContext, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testContext),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name    string
		resp    *modelexperiments.CreateInstanceTokenResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name:    "empty response",
			resp:    &modelexperiments.CreateInstanceTokenResponse{},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := fixtureInputModel(func(model *inputModel) {
				model.OutputFormat = ""
				model.Async = false
			})
			if err := outputResult(params.Printer, model, "test-project", "test-instance", tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
