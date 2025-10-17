package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	testRegion         = "eu02"
	testCredentialsRef = "credentials-test"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &loadbalancer.APIClient{}
	testProjectId = uuid.NewString()
)

type loadBalancerClientMocked struct {
	getCredentialsError    bool
	getCredentialsResponse *loadbalancer.GetCredentialsResponse
}

func (c *loadBalancerClientMocked) UpdateCredentials(ctx context.Context, projectId, region, credentialsRef string) loadbalancer.ApiUpdateCredentialsRequest {
	return testClient.UpdateCredentials(ctx, projectId, region, credentialsRef)
}

func (c *loadBalancerClientMocked) GetCredentialsExecute(_ context.Context, _, _, _ string) (*loadbalancer.GetCredentialsResponse, error) {
	if c.getCredentialsError {
		return nil, fmt.Errorf("get credentials failed")
	}
	return c.getCredentialsResponse, nil
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testCredentialsRef,
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
		displayNameFlag:           "name",
		usernameFlag:              "username",
		passwordFlag:              "pwd",
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
		DisplayName:    utils.Ptr("name"),
		Username:       utils.Ptr("username"),
		Password:       utils.Ptr("pwd"),
		CredentialsRef: testCredentialsRef,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiUpdateCredentialsRequest)) loadbalancer.ApiUpdateCredentialsRequest {
	request := testClient.UpdateCredentials(testCtx, testProjectId, testRegion, testCredentialsRef)
	request = request.UpdateCredentialsPayload(loadbalancer.UpdateCredentialsPayload{
		DisplayName: utils.Ptr("name"),
		Username:    utils.Ptr("username"),
		Password:    utils.Ptr("pwd"),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureGetCredentialsResponse(mods ...func(response *loadbalancer.GetCredentialsResponse)) *loadbalancer.GetCredentialsResponse {
	response := &loadbalancer.GetCredentialsResponse{
		Credential: &loadbalancer.CredentialsResponse{
			DisplayName: utils.Ptr("name"),
			Username:    utils.Ptr("username"),
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
			description: "credentials ref invalid 1",
			argValues:   []string{""},
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
	tests := []struct {
		description            string
		model                  *inputModel
		expectedRequest        loadbalancer.ApiUpdateCredentialsRequest
		getCredentialsFails    bool
		getCredentialsResponse *loadbalancer.GetCredentialsResponse
		isValid                bool
	}{
		{
			description:            "base",
			model:                  fixtureInputModel(),
			expectedRequest:        fixtureRequest(),
			getCredentialsResponse: fixtureGetCredentialsResponse(),
			isValid:                true,
		},
		{
			description: "no display name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.DisplayName = nil
				},
			),
			expectedRequest:        fixtureRequest(),
			getCredentialsResponse: fixtureGetCredentialsResponse(),
			isValid:                true,
		},
		{
			description: "no username name",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.Username = nil
				},
			),
			expectedRequest:        fixtureRequest(),
			getCredentialsResponse: fixtureGetCredentialsResponse(),
			isValid:                true,
		},
		{
			description:         "get credentials fails",
			model:               fixtureInputModel(),
			getCredentialsFails: true,
			isValid:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &loadBalancerClientMocked{
				getCredentialsError:    tt.getCredentialsFails,
				getCredentialsResponse: tt.getCredentialsResponse,
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
