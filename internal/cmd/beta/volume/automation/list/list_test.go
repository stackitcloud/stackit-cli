package list

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &automation.APIClient{DefaultAPI: &automation.DefaultAPIService{}}
	testProjectId = uuid.NewString()
)

const (
	testLimit  int64 = 10
	testRegion       = "eu01"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		limitFlag:                 strconv.FormatInt(testLimit, 10),
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
		Limit: new(testLimit),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *automation.ApiListVolumeAutomationsRequest)) automation.ApiListVolumeAutomationsRequest {
	request := testClient.DefaultAPI.ListVolumeAutomations(testCtx, testProjectId, testRegion)
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
			description: "no limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, limitFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = nil
			}),
		},
		{
			description:   "no values",
			flagValues:    map[string]string{},
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "missing project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "invalid limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
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
		expectedRequest automation.ApiListVolumeAutomationsRequest
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

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest, automation.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match (-want, +got): %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		items        []automation.ListAutomationsItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty list",
			args: args{
				projectLabel: testProjectId,
				items:        []automation.ListAutomationsItem{},
			},
			wantErr: false,
		},
		{
			name: "items list",
			args: args{
				projectLabel: testProjectId,
				items: []automation.ListAutomationsItem{
					{
						Id:          uuid.NewString(),
						Name:        new("my-automation"),
						Description: new("my-description"),
						TemplateId:  new(uuid.NewString()),
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
