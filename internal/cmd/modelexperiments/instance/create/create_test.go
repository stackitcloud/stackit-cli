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

type testCtxKey struct{}

const (
	testRegion = "eu01"
	testName   = "example"
)

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &modelexperiments.APIClient{DefaultAPI: modelexperiments.DefaultAPIServiceMock{}}
	testProjectId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		nameFlag:                  testName,
		globalflags.RegionFlag:    testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		Name: testName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *modelexperiments.ApiCreateInstanceRequest)) modelexperiments.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(modelexperiments.CreateInstancePayload{
		Name: testName,
	})
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
			description: "with description and retention",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[descriptionFlag] = "team tracking server"
				flagValues[retentionFlag] = "30d"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("team tracking server")
				model.Retention = utils.Ptr("30d")
			}),
		},
		{
			description: "with labels",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "env=prod,team=ml"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{"env": "prod", "team": "ml"}
			}),
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
		expectedRequest modelexperiments.ApiCreateInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with description",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("team tracking server")
			}),
			expectedRequest: fixtureRequest(func(request *modelexperiments.ApiCreateInstanceRequest) {
				payload := modelexperiments.CreateInstancePayload{
					Name:        testName,
					Description: utils.Ptr("team tracking server"),
				}
				*request = request.CreateInstancePayload(payload)
			}),
		},
		{
			description: "with labels",
			model: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{"env": "prod", "team": "ml"}
			}),
			expectedRequest: fixtureRequest(func(request *modelexperiments.ApiCreateInstanceRequest) {
				payload := modelexperiments.CreateInstancePayload{
					Name:   testName,
					Labels: &map[string]string{"env": "prod", "team": "ml"},
				}
				*request = request.CreateInstancePayload(payload)
			}),
		},
		{
			description: "with retention",
			model: fixtureInputModel(func(model *inputModel) {
				model.Retention = utils.Ptr("30d")
			}),
			expectedRequest: fixtureRequest(func(request *modelexperiments.ApiCreateInstanceRequest) {
				payload := modelexperiments.CreateInstancePayload{
					Name:                       testName,
					DeletedExperimentRetention: utils.Ptr("30d"),
				}
				*request = request.CreateInstancePayload(payload)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
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
		resp    *modelexperiments.CreateInstanceResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name:    "empty response",
			resp:    &modelexperiments.CreateInstanceResponse{},
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
			if err := outputResult(params.Printer, model, "test-project", tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
