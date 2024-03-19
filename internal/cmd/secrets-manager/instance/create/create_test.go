package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	testACL1 = "1.2.3.4/24"
	testACL2 = "4.3.2.1/12"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &secretsmanager.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:    testProjectId,
		instanceNameFlag: "example",
		aclFlag:          testACL1,
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
		},
		InstanceName: utils.Ptr("example"),
		Acls:         utils.Ptr([]string{testACL1}),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *secretsmanager.ApiCreateInstanceRequest)) secretsmanager.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId)
	request = request.CreateInstancePayload(secretsmanager.CreateInstancePayload{
		Name: utils.Ptr("example"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureUpdateACLsRequest(mods ...func(request *secretsmanager.ApiUpdateACLsRequest)) secretsmanager.ApiUpdateACLsRequest {
	request := testClient.UpdateACLs(testCtx, testProjectId, testInstanceId)
	request = request.UpdateACLsPayload(secretsmanager.UpdateACLsPayload{
		Cidrs: utils.Ptr([]secretsmanager.AclUpdate{
			{Cidr: utils.Ptr(testACL1)},
		})})

	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "zero values",
			flagValues: map[string]string{
				projectIdFlag:    testProjectId,
				instanceNameFlag: "",
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				InstanceName: utils.Ptr(""),
			},
		},
		{
			description: "instance name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceNameFlag)
			}),
			isValid: false,
		},
		{
			description: "acl missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = nil
			}),
		},
		{
			description: "acl empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "multiple acls",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = testACL1 + "," + testACL2
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				*model.Acls = append(*model.Acls, testACL2)
			}),
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
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd()
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

			model, err := parseInput(cmd)
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

func TestBuildCreateInstanceRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest secretsmanager.ApiCreateInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildCreateInstanceRequest(testCtx, tt.model, testClient)

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
func TestBuildCreateACLRequests(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest secretsmanager.ApiUpdateACLsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureUpdateACLsRequest(),
		},
		{
			description: "multiple ACLs",
			model: fixtureInputModel(func(model *inputModel) {
				*model.Acls = append(*model.Acls, testACL2)
			}),
			expectedRequest: fixtureUpdateACLsRequest().UpdateACLsPayload(secretsmanager.UpdateACLsPayload{
				Cidrs: utils.Ptr([]secretsmanager.AclUpdate{
					{Cidr: utils.Ptr(testACL1)},
					{Cidr: utils.Ptr(testACL2)},
				})}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildUpdateACLsRequest(testCtx, tt.model, testInstanceId, testClient)

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
