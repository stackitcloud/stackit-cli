package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	observability "github.com/stackitcloud/stackit-sdk-go/services/observability/v1api"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &observability.APIClient{DefaultAPI: &observability.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:  testProjectId,
		instanceIdFlag: testInstanceId,
		aclFlag:        "198.51.100.14/24",
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
		Acls:       []string{"198.51.100.14/24"},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *observability.ApiUpdateACLRequest)) observability.ApiUpdateACLRequest {
	request := testClient.DefaultAPI.UpdateACL(testCtx, testInstanceId, testProjectId)
	request = request.UpdateACLPayload(observability.UpdateACLPayload{
		Acl: []string{"198.51.100.14/24"},
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
		aclValues     []string
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
			description: "multiple acls",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = "198.51.100.14/24,1.2.3.4/32"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = []string{"198.51.100.14/24", "1.2.3.4/32"}
			}),
		},
		{
			description: "repeated acl flags",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			aclValues: []string{"198.51.100.14/24", "1.2.3.4/32"},
			isValid:   true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = []string{"198.51.100.14/24", "1.2.3.4/32"}
			}),
		},
		{
			description: "empty acl removes all acls",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acls = []string{}
			}),
		},
		{
			description: "acl missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			isValid: false,
		},
		{
			description: "acl invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = "invalid-cidr"
			}),
			isValid: false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
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
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				aclFlag: tt.aclValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest observability.ApiUpdateACLRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "multiple acls",
			model: fixtureInputModel(func(model *inputModel) {
				model.Acls = []string{"198.51.100.14/24", "1.2.3.4/32"}
			}),
			expectedRequest: fixtureRequest(func(request *observability.ApiUpdateACLRequest) {
				*request = request.UpdateACLPayload(observability.UpdateACLPayload{
					Acl: []string{"198.51.100.14/24", "1.2.3.4/32"},
				})
			}),
		},
		{
			description: "empty acls",
			model: fixtureInputModel(func(model *inputModel) {
				model.Acls = []string{}
			}),
			expectedRequest: fixtureRequest(func(request *observability.ApiUpdateACLRequest) {
				*request = request.UpdateACLPayload(observability.UpdateACLPayload{
					Acl: []string{},
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, observability.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
