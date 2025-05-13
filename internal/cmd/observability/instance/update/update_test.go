package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	listPlansError      bool
	listPlansResponse   *observability.PlansResponse
	getInstanceError    bool
	getInstanceResponse *observability.GetInstanceResponse
}

func (c *observabilityClientMocked) UpdateInstance(ctx context.Context, instanceId, projectId string) observability.ApiUpdateInstanceRequest {
	return testClient.UpdateInstance(ctx, instanceId, projectId)
}

func (c *observabilityClientMocked) ListPlansExecute(_ context.Context, _ string) (*observability.PlansResponse, error) {
	if c.listPlansError {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listPlansResponse, nil
}

func (c *observabilityClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*observability.GetInstanceResponse, error) {
	if c.getInstanceError {
		return nil, fmt.Errorf("get instance failed")
	}
	return c.getInstanceResponse, nil
}

const (
	testInstanceName = "example-instance-name"
)

var (
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testPlanId     = uuid.NewString()
	testNewPlanId  = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		planIdFlag:    testNewPlanId,
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
		InstanceId: testInstanceId,
		PlanId:     utils.Ptr(testNewPlanId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *observability.ApiUpdateInstanceRequest)) observability.ApiUpdateInstanceRequest {
	request := testClient.UpdateInstance(testCtx, testInstanceId, testProjectId)
	request = request.UpdateInstancePayload(observability.UpdateInstancePayload{
		PlanId: utils.Ptr(testNewPlanId),
		Name:   utils.Ptr(testInstanceName),
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
				Id:   utils.Ptr(testNewPlanId),
			},
		},
	}
	for _, mod := range mods {
		mod(response)
	}
	return response
}

func fixtureGetInstanceResponse(mods ...func(response *observability.GetInstanceResponse)) *observability.GetInstanceResponse {
	response := &observability.GetInstanceResponse{
		PlanId: utils.Ptr(testPlanId),
		Name:   utils.Ptr(testInstanceName),
	}
	for _, mod := range mods {
		mod(response)
	}
	return response
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
			description: "with plan name",
			argValues:   fixtureArgValues(),
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
			description: "with new instance name",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceNameFlag] = "new-instance-name"
				delete(flagValues, planIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlanId = nil
				model.PlanName = ""
				model.InstanceName = utils.Ptr("new-instance-name")
			}),
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
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
			cmd := NewCmd(&params.CmdParams{Printer: p})
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
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
		description         string
		model               *inputModel
		expectedRequest     observability.ApiUpdateInstanceRequest
		getPlansFails       bool
		getPlansResponse    *observability.PlansResponse
		getInstanceFails    bool
		getInstanceResponse *observability.GetInstanceResponse
		isValid             bool
	}{
		{
			description:         "base",
			model:               fixtureInputModel(),
			expectedRequest:     fixtureRequest(),
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			isValid:             true,
		},
		{
			description: "use plan name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = "example-plan-name"
				},
			),
			expectedRequest:     fixtureRequest(),
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			isValid:             true,
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
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			isValid:             false,
		},
		{
			description: "plan id not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = utils.Ptr(uuid.NewString())
				},
			),
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			isValid:             false,
		},
		{
			description: "plan id, no instance name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.InstanceName = nil
				},
			),
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			expectedRequest:     fixtureRequest(),
			isValid:             true,
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
			getPlansResponse:    fixturePlansResponse(),
			getInstanceResponse: fixtureGetInstanceResponse(),
			expectedRequest:     fixtureRequest(),
			isValid:             true,
		},
		{
			description: "instance name, no plan info",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = ""
					model.InstanceName = utils.Ptr("new-instance-name")
				},
			),
			getInstanceResponse: fixtureGetInstanceResponse(),
			expectedRequest: fixtureRequest().
				UpdateInstancePayload(observability.UpdateInstancePayload{
					PlanId: utils.Ptr(testPlanId),
					Name:   utils.Ptr("new-instance-name"),
				}),
			isValid: true,
		},
		{
			description: "instance name, no plan info, get instance fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.PlanId = nil
					model.PlanName = ""
					model.InstanceName = utils.Ptr("new-instance-name")
				},
			),
			getInstanceFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &observabilityClientMocked{
				listPlansError:      tt.getPlansFails,
				listPlansResponse:   tt.getPlansResponse,
				getInstanceError:    tt.getInstanceFails,
				getInstanceResponse: tt.getInstanceResponse,
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
