package set

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const (
	testRegion   = "eu01"
	testUserName = "jane.doe"
)

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &scf.APIClient{
		DefaultAPI: &scf.DefaultAPIService{},
	}
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
	testSpaceId        = uuid.NewString()
	testUserGuid       = uuid.NewString()
)

func fixtureArgValues() []string {
	return []string{testSpaceId}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		organizationIdFlag:        testOrganizationId,
		roleFlag.Name():           string(scf.SPACEROLETYPE_SPACE_DEVELOPER),
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
		SpaceId:        testSpaceId,
		Role:           scf.SPACEROLETYPE_SPACE_DEVELOPER,
		UserGuid:       new(testUserGuid),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "user name instead of guid",
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
			description: "user guid and name both set",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[userNameFlag] = testUserName
			}),
			isValid: false,
		},
		{
			description: "no user given",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, userGuidFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid role",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[roleFlag.Name()] = "organization_manager"
			}),
			isValid: false,
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
			description: "user guid invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[userGuidFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "organization id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "space id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	expected := testClient.DefaultAPI.CreateSpaceRole(testCtx, testProjectId, testRegion, testOrganizationId, testSpaceId).
		CreateSpaceRolePayload(scf.CreateSpaceRolePayload{
			Type:     scf.SPACEROLETYPE_SPACE_DEVELOPER,
			UserGuid: new(testUserGuid),
		})
	request := buildRequest(testCtx, fixtureInputModel(), testClient.DefaultAPI)

	diff := cmp.Diff(request, expected,
		cmp.AllowUnexported(expected),
		cmpopts.EquateComparable(testCtx),
		cmpopts.EquateComparable(testClient.DefaultAPI),
	)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name    string
		model   *inputModel
		resp    *scf.SpaceRoleCreateResponse
		wantErr bool
	}{
		{
			name:    "default output",
			model:   fixtureInputModel(),
			resp:    &scf.SpaceRoleCreateResponse{},
			wantErr: false,
		},
		{
			name:    "json output",
			model:   fixtureInputModel(),
			resp:    &scf.SpaceRoleCreateResponse{},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.model, "space", tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
