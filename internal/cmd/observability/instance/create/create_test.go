package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &observability.APIClient{}

type observabilityClientMocked struct {
	returnError       bool
	listPlansResponse *observability.PlansResponse
}

func (c *observabilityClientMocked) CreateInstance(ctx context.Context, projectId string) observability.ApiCreateInstanceRequest {
	return testClient.CreateInstance(ctx, projectId)
}

func (c *observabilityClientMocked) ListPlansExecute(_ context.Context, _ string) (*observability.PlansResponse, error) {
	if c.returnError {
		return nil, fmt.Errorf("list plans failed")
	}
	return c.listPlansResponse, nil
}

var testProjectId = uuid.NewString()
var testPlanId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:    testProjectId,
		instanceNameFlag: "example-name",
		planIdFlag:       testPlanId,
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
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceName: utils.Ptr("example-name"),
		PlanId:       utils.Ptr(testPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *observability.ApiCreateInstanceRequest)) observability.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId)
	request = request.CreateInstancePayload(observability.CreateInstancePayload{
		Name:   utils.Ptr("example-name"),
		PlanId: utils.Ptr(testPlanId),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePlansResponse(mods ...func(response *observability.PlansResponse)) *observability.PlansResponse {
	response := &observability.PlansResponse{
		Plans: &[]observability.Plan{
			{
				Name: utils.Ptr("example-plan-name"),
				Id:   utils.Ptr(testPlanId),
			},
		},
	}
	for _, mod := range mods {
		mod(response)
	}
	return response
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
			description: "with plan name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = "plan-name"
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:    testProjectId,
				planIdFlag:       testPlanId,
				instanceNameFlag: "",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				PlanId:       utils.Ptr(testPlanId),
				InstanceName: utils.Ptr(""),
			},
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "invalid with plan ID and plan name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[planNameFlag] = "plan-name"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description      string
		model            *inputModel
		expectedRequest  observability.ApiCreateInstanceRequest
		getPlansFails    bool
		getPlansResponse *observability.PlansResponse
		isValid          bool
	}{
		{
			description:      "base",
			model:            fixtureInputModel(),
			expectedRequest:  fixtureRequest(),
			getPlansResponse: fixturePlansResponse(),
			isValid:          true,
		},
		{
			description: "use plan name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
				},
			),
			expectedRequest:  fixtureRequest(),
			getPlansResponse: fixturePlansResponse(),
			isValid:          true,
		},
		{
			description: "get plans fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
				},
			),
			getPlansFails: true,
			isValid:       false,
		},
		{
			description: "plan name not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "non-existent-plan"
				},
			),
			getPlansResponse: fixturePlansResponse(),
			isValid:          false,
		},
		{
			description: "plan id not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = utils.Ptr(uuid.NewString())
				},
			),
			getPlansResponse: fixturePlansResponse(),
			isValid:          false,
		},
		{
			description: "plan id, no instance name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.InstanceName = nil
				},
			),
			getPlansResponse: fixturePlansResponse(),
			expectedRequest: testClient.CreateInstance(testCtx, testProjectId).
				CreateInstancePayload(observability.CreateInstancePayload{PlanId: utils.Ptr(testPlanId)}),
			isValid: true,
		},
		{
			description: "plan name, no instance name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
					model.InstanceName = nil
				},
			),
			getPlansResponse: fixturePlansResponse(),
			expectedRequest: testClient.CreateInstance(testCtx, testProjectId).
				CreateInstancePayload(observability.CreateInstancePayload{PlanId: utils.Ptr(testPlanId)}),
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &observabilityClientMocked{
				returnError:       tt.getPlansFails,
				listPlansResponse: tt.getPlansResponse,
			}
			request, err := buildRequest(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			if !tt.isValid {
				t.Fatal("expected error but none thrown")
			}

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
