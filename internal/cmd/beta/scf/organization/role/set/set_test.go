package set

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
)

var (
	testCtx            = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient         = &scf.DefaultAPIService{}
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
	testUserGuid       = uuid.NewString()
	testUserName       = "john.stackit"
	testRole           = string(scf.ORGROLETYPE_ORGANIZATION_MANAGER)
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testOrganizationId,
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
		roleFlag.Name():           testRole,
		userGuidFlag:              testUserGuid,
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
		OrganizationId: testOrganizationId,
		Role:           scf.OrgRoleType(testRole),
		UserGuid:       new(testUserGuid),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *scf.ApiCreateOrgRoleRequest)) scf.ApiCreateOrgRoleRequest {
	request := testClient.CreateOrgRole(testCtx, testProjectId, testRegion, testOrganizationId).
		CreateOrgRolePayload(scf.CreateOrgRolePayload{
			Type:     scf.OrgRoleType(testRole),
			UserGuid: new(testUserGuid),
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
			description:   "base with user guid",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "with user name",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, userGuidFlag)
				flagValues[userNameFlag] = testUserName
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.UserGuid = nil
				model.UserName = new(testUserName)
			}),
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
			description: "organization id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "role missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, roleFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "role invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[roleFlag.Name()] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "user guid and user name given",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[userNameFlag] = testUserName
			}),
			isValid: false,
		},
		{
			description: "neither user guid nor user name given",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, userGuidFlag)
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
		expectedRequest scf.ApiCreateOrgRoleRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with user name",
			model: fixtureInputModel(func(model *inputModel) {
				model.UserGuid = nil
				model.UserName = new(testUserName)
			}),
			expectedRequest: fixtureRequest(func(request *scf.ApiCreateOrgRoleRequest) {
				*request = request.CreateOrgRolePayload(scf.CreateOrgRolePayload{
					Type:     scf.OrgRoleType(testRole),
					UserName: new(testUserName),
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.EquateComparable(testClient),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name         string
		model        *inputModel
		resp         *scf.OrgRoleResponse
		outputFormat string
		wantErr      bool
	}{
		{
			name:    "default output",
			model:   fixtureInputModel(),
			resp:    fixtureOrgRoleResponse(),
			wantErr: false,
		},
		{
			name:         "json output",
			model:        fixtureInputModel(),
			resp:         fixtureOrgRoleResponse(),
			outputFormat: print.JSONOutputFormat,
			wantErr:      false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.model, tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func fixtureOrgRoleResponse() *scf.OrgRoleResponse {
	return &scf.OrgRoleResponse{
		Guid:       uuid.NewString(),
		OrgId:      testOrganizationId,
		PlatformId: uuid.NewString(),
		ProjectId:  testProjectId,
		Region:     testRegion,
		Type:       scf.OrgRoleType(testRole),
	}
}
