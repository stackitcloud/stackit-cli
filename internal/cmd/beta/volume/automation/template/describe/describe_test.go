package describe

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &automation.APIClient{DefaultAPI: &automation.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testTemplateId = uuid.NewString()
)

const testRegion = "eu01"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testTemplateId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

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
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		TemplateId: testTemplateId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *automation.ApiGetVolumeTemplateRequest)) automation.ApiGetVolumeTemplateRequest {
	request := testClient.DefaultAPI.GetVolumeTemplate(testCtx, testProjectId, testRegion, testTemplateId)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description:   "no values",
			argValues:     []string{},
			flagValues:    map[string]string{},
			isValid:       false,
			expectedModel: nil,
		},
		{
			description:   "missing args",
			argValues:     []string{},
			flagValues:    fixtureFlagValues(),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "missing project id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid:       false,
			expectedModel: nil,
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
		isValid         bool
		expectedRequest automation.ApiGetVolumeTemplateRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, automation.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		templateId   string
		projectLabel string
		template     *automation.GetVolumeTemplateResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "empty response",
			args: args{
				templateId:   testTemplateId,
				projectLabel: testProjectId,
				template:     &automation.GetVolumeTemplateResponse{},
			},
			wantErr: false,
		},
		{
			name: "empty output field",
			args: args{
				templateId:   testTemplateId,
				projectLabel: testProjectId,
				template: &automation.GetVolumeTemplateResponse{
					Output: &automation.VolumeOutput{},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.templateId, tt.args.projectLabel, tt.args.template); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
