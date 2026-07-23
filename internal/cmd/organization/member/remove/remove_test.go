package remove

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	authorization "github.com/stackitcloud/stackit-sdk-go/services/authorization/v2api"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &authorization.APIClient{DefaultAPI: &authorization.DefaultAPIService{}}
var testOrganizationID = "some-organization-id"
var testSubject = "someone@domain.com"
var testRole = "reader"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testSubject,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		organizationIdFlag: testOrganizationID,
		roleFlag:           testRole,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
		OrganizationId:  testOrganizationID,
		Subject:         testSubject,
		Role:            testRole,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *authorization.ApiRemoveMembersRequest)) authorization.ApiRemoveMembersRequest {
	request := testClient.DefaultAPI.RemoveMembers(testCtx, testOrganizationID)
	request = request.RemoveMembersPayload(authorization.RemoveMembersPayload{
		Members: []authorization.Member{
			{
				Subject: testSubject,
				Role:    testRole,
			},
		},
		ResourceType: organizationResourceType,
		ForceRemove:  utils.Ptr(false),
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
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "with force",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[forceFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.Force = true
				}),
		},
		{
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no flags",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "organization id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "role missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, roleFlag)
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
		expectedRequest authorization.ApiRemoveMembersRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with force",
			model: fixtureInputModel(func(model *inputModel) {
				model.Force = true
			}),
			expectedRequest: testClient.DefaultAPI.RemoveMembers(testCtx, testOrganizationID).
				RemoveMembersPayload(authorization.RemoveMembersPayload{
					Members: []authorization.Member{
						{
							Subject: testSubject,
							Role:    testRole,
						},
					},
					ResourceType: organizationResourceType,
					ForceRemove:  utils.Ptr(true),
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, authorization.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
