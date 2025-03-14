package delete

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
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/wait"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &postgresflex.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testRegion = "eu01"

type postgresFlexClientMocked struct {
	getInstanceFails bool
	getInstanceResp  *postgresflex.InstanceResponse
}

func (c *postgresFlexClientMocked) GetInstanceExecute(_ context.Context, _, _, _ string) (*postgresflex.InstanceResponse, error) {
	if c.getInstanceFails {
		return nil, fmt.Errorf("get instance failed")
	}
	return c.getInstanceResp, nil
}

func (c *postgresFlexClientMocked) ListVersionsExecute(_ context.Context, _, _ string) (*postgresflex.ListVersionsResponse, error) {
	// Not used in testing
	return nil, nil
}
func (c *postgresFlexClientMocked) GetUserExecute(_ context.Context, _, _, _, _ string) (*postgresflex.GetUserResponse, error) {
	// Not used in testing
	return nil, nil
}

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
		InstanceId: testInstanceId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureDeleteRequest(mods ...func(request *postgresflex.ApiDeleteInstanceRequest)) postgresflex.ApiDeleteInstanceRequest {
	request := testClient.DeleteInstance(testCtx, testProjectId, testRegion, testInstanceId)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureForceDeleteRequest(mods ...func(request *postgresflex.ApiForceDeleteInstanceRequest)) postgresflex.ApiForceDeleteInstanceRequest {
	request := testClient.ForceDeleteInstance(testCtx, testProjectId, testRegion, testInstanceId)
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
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
				t.Fatalf("error parsing input: %v", err)
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

func TestBuildDeleteRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest postgresflex.ApiDeleteInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureDeleteRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildDeleteRequest(testCtx, tt.model, testClient)

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

func TestBuildForceDeleteRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest postgresflex.ApiForceDeleteInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureForceDeleteRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildForceDeleteRequest(testCtx, tt.model, testClient)

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

func TestCheckIfInstanceIsDeleted(t *testing.T) {
	tests := []struct {
		description           string
		model                 *inputModel
		expectedToDelete      bool
		expectedToForceDelete bool
		getInstanceResponse   *postgresflex.InstanceResponse
		getInstanceFails      bool
		isValid               bool
	}{
		{
			description:           "delete instance state Ready",
			model:                 fixtureInputModel(),
			expectedToDelete:      true,
			expectedToForceDelete: false,
			getInstanceResponse: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Status: utils.Ptr(wait.InstanceStateSuccess),
				},
			},
			isValid: true,
		},
		{
			description: "force delete instance state Ready",
			model: fixtureInputModel(func(model *inputModel) {
				model.ForceDelete = true
			}),
			expectedToDelete:      true,
			expectedToForceDelete: true,
			getInstanceResponse: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Status: utils.Ptr(wait.InstanceStateSuccess),
				},
			},
			isValid: true,
		},
		{
			description: "force delete instance state Deleted",
			model: fixtureInputModel(func(model *inputModel) {
				model.ForceDelete = true
			}),
			expectedToDelete:      false,
			expectedToForceDelete: true,
			getInstanceResponse: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Status: utils.Ptr(wait.InstanceStateDeleted),
				},
			},
			isValid: true,
		},
		{
			description: "delete instance state Deleted",
			model:       fixtureInputModel(),
			getInstanceResponse: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Status: utils.Ptr(wait.InstanceStateDeleted),
				},
			},
			isValid: false,
		},
		{
			description:      "delete instance get instance fails",
			model:            fixtureInputModel(),
			getInstanceFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &postgresFlexClientMocked{
				getInstanceResp:  tt.getInstanceResponse,
				getInstanceFails: tt.getInstanceFails,
			}

			toDelete, toForceDelete, err := getNextOperations(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error checking if instance is deleted: %v", err)
			}

			if toDelete != tt.expectedToDelete {
				t.Fatalf("toDelete does not match: got %v, expected %v", toDelete, tt.expectedToDelete)
			}

			if toForceDelete != tt.expectedToForceDelete {
				t.Fatalf("toForceDelete does not match: got %v, expected %v", toForceDelete, tt.expectedToForceDelete)
			}
		})
	}
}
