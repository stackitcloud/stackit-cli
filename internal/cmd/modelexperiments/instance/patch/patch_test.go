package patch

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

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

var testClient = &modelexperiments.APIClient{
	DefaultAPI: modelexperiments.DefaultAPIServiceMock{},
}

var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(
	mods ...func(*map[string]string),
) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: "eu01",
	}

	for _, mod := range mods {
		mod(&flagValues)
	}

	return flagValues
}

func fixtureInputModel(
	mods ...func(*inputModel),
) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    "eu01",
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId: testInstanceId,
	}

	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(
	mods ...func(*modelexperiments.ApiPartialUpdateInstanceRequest),
) modelexperiments.ApiPartialUpdateInstanceRequest {
	request := testClient.DefaultAPI.PartialUpdateInstance(
		testCtx,
		testProjectId,
		"eu01",
		testInstanceId,
	)

	payload := modelexperiments.PartialUpdateInstancePayload{
		Name: utils.Ptr("example"),
	}

	request = request.PartialUpdateInstancePayload(payload)

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
			description: "name",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues *map[string]string) {
				(*flagValues)[nameFlag] = "example"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = utils.Ptr("example")
			}),
		},
		{
			description: "description",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues *map[string]string) {
				(*flagValues)[descriptionFlag] = "team tracking server"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("team tracking server")
			}),
		},
		{
			description: "retention",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues *map[string]string) {
				(*flagValues)[retentionFlag] = "30d"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Retention = utils.Ptr("30d")
			}),
		},
		{
			description: "no update flags",
			argValues:   []string{testInstanceId},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   []string{testInstanceId},
			flagValues: fixtureFlagValues(func(flagValues *map[string]string) {
				delete(*flagValues, projectIdFlag)
				(*flagValues)[nameFlag] = "example"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(
				t,
				NewCmd,
				parseInput,
				tt.expectedModel,
				tt.argValues,
				tt.flagValues,
				tt.isValid,
			)
		})
	}
}

func TestBuildPatchInstanceRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest modelexperiments.ApiPartialUpdateInstanceRequest
	}{
		{
			description: "name",
			model: fixtureInputModel(func(model *inputModel) {
				model.Name = utils.Ptr("example")
			}),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "description",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = utils.Ptr("team tracking server")
			}),
			expectedRequest: func() modelexperiments.ApiPartialUpdateInstanceRequest {
				request := testClient.DefaultAPI.PartialUpdateInstance(
					testCtx,
					testProjectId,
					"eu01",
					testInstanceId,
				)

				request = request.PartialUpdateInstancePayload(
					modelexperiments.PartialUpdateInstancePayload{
						Description: utils.Ptr("team tracking server"),
					},
				)

				return request
			}(),
		},
		{
			description: "retention",
			model: fixtureInputModel(func(model *inputModel) {
				model.Retention = utils.Ptr("30d")
			}),
			expectedRequest: func() modelexperiments.ApiPartialUpdateInstanceRequest {
				request := testClient.DefaultAPI.PartialUpdateInstance(
					testCtx,
					testProjectId,
					"eu01",
					testInstanceId,
				)

				request = request.PartialUpdateInstancePayload(
					modelexperiments.PartialUpdateInstancePayload{
						DeletedExperimentRetention: utils.Ptr("30d"),
					},
				)

				return request
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildPatchInstanceRequest(
				testCtx,
				tt.model,
				testClient,
			)

			diff := cmp.Diff(
				request,
				tt.expectedRequest,
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
		resp    *modelexperiments.PartialUpdateInstanceResponse
		wantErr bool
	}{
		{
			name:    "nil response",
			resp:    nil,
			wantErr: true,
		},
		{
			name: "empty instance",
			resp: &modelexperiments.PartialUpdateInstanceResponse{
				Instance: modelexperiments.Instance{},
			},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(
				params.Printer,
				"",
				tt.resp,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"outputResult() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
