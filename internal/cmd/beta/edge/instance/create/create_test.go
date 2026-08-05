package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testProjectId  = uuid.NewString()
	testPlanId     = uuid.NewString()
	testInstanceId = uuid.NewString()
	testClient     = &edge.APIClient{DefaultAPI: &edge.DefaultAPIService{}}
)

const (
	testRegion      = "eu01"
	testName        = "test"
	testDescription = "Initial instance description"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testName,
		descriptionFlag:           testDescription,
		planIdFlag:                testPlanId,
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
		DisplayName: testName,
		Description: utils.Ptr(testDescription),
		PlanId:      testPlanId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *edge.ApiCreateInstanceRequest)) edge.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *edge.CreateInstancePayload)) edge.CreateInstancePayload {
	payload := edge.CreateInstancePayload{
		DisplayName: testName,
		Description: utils.Ptr(testDescription),
		PlanId:      testPlanId,
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "plan id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, planIdFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testUtils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest edge.ApiCreateInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Region:    testRegion,
					Verbosity: globalflags.VerbosityDefault,
				},
				DisplayName: testName,
				PlanId:      testPlanId,
			},
			expectedRequest: testClient.DefaultAPI.
				CreateInstance(testCtx, testProjectId, testRegion).
				CreateInstancePayload(edge.CreateInstancePayload{
					DisplayName: testName,
					PlanId:      testPlanId,
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, edge.DefaultAPIService{}),
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
		model        *inputModel
		instance     *edge.Instance
		projectLabel string
	}

	tests := []struct {
		description string
		wantErr     bool
		args        args
	}{
		{
			description: "no instance",
			wantErr:     true,
			args: args{
				model: fixtureInputModel(),
			},
		},
		{
			description: "output json",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				instance: &edge.Instance{},
			},
		},
		{
			description: "output yaml",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.YAMLOutputFormat
				}),
				instance: &edge.Instance{},
			},
		},
		{
			description: "output default",
			args: args{
				model:    fixtureInputModel(),
				instance: &edge.Instance{Id: testInstanceId},
			},
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.instance); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
