package enable

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
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
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

type observabilityClientMocked struct {
	getGrafanaConfigsFails bool
	getGrafanaConfigsResp  *observability.GrafanaConfigs
}

func (c *observabilityClientMocked) GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*observability.GetInstanceResponse, error) {
	return testClient.GetInstanceExecute(ctx, instanceId, projectId)
}

func (c *observabilityClientMocked) UpdateGrafanaConfigs(ctx context.Context, instanceId, projectId string) observability.ApiUpdateGrafanaConfigsRequest {
	return testClient.UpdateGrafanaConfigs(ctx, instanceId, projectId)
}

func (c *observabilityClientMocked) GetGrafanaConfigsExecute(_ context.Context, _, _ string) (*observability.GrafanaConfigs, error) {
	if c.getGrafanaConfigsFails {
		return nil, fmt.Errorf("get payload failed")
	}
	return c.getGrafanaConfigsResp, nil
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
		projectIdFlag: testProjectId,
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
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureGrafanaConfigs(mods ...func(gc *observability.GrafanaConfigs)) *observability.GrafanaConfigs {
	gc := observability.GrafanaConfigs{
		GenericOauth: &observability.GrafanaOauth{
			ApiUrl:              utils.Ptr("apiUrl"),
			AuthUrl:             utils.Ptr("authUrl"),
			Enabled:             utils.Ptr(true),
			Name:                utils.Ptr("name"),
			OauthClientId:       utils.Ptr("oauthClientId"),
			OauthClientSecret:   utils.Ptr("oauthClientSecret"),
			RoleAttributePath:   utils.Ptr("roleAttributePath"),
			RoleAttributeStrict: utils.Ptr(true),
			Scopes:              utils.Ptr("scopes"),
			TokenUrl:            utils.Ptr("tokenUrl"),
			UsePkce:             utils.Ptr(true),
		},
		PublicReadAccess: utils.Ptr(false),
		UseStackitSso:    utils.Ptr(false),
	}
	for _, mod := range mods {
		mod(&gc)
	}
	return &gc
}

func fixturePayload(mods ...func(payload *observability.UpdateGrafanaConfigsPayload)) *observability.UpdateGrafanaConfigsPayload {
	payload := &observability.UpdateGrafanaConfigsPayload{
		GenericOauth:     observabilityUtils.ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
		PublicReadAccess: fixtureGrafanaConfigs().PublicReadAccess,
		UseStackitSso:    utils.Ptr(true),
	}
	for _, mod := range mods {
		mod(payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *observability.ApiUpdateGrafanaConfigsRequest)) observability.ApiUpdateGrafanaConfigsRequest {
	request := testClient.UpdateGrafanaConfigs(testCtx, testInstanceId, testProjectId)
	request = request.UpdateGrafanaConfigsPayload(*fixturePayload())
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
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
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
		getGrafanaConfigsFails bool
		getGrafanaConfigsResp  *observability.GrafanaConfigs
		isValid                bool
		expectedRequest        observability.ApiUpdateGrafanaConfigsRequest
	}{
		{
			description:           "base",
			model:                 fixtureInputModel(),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(),
			isValid:               true,
			expectedRequest:       fixtureRequest(),
		},
		{
			description: "nil generic oauth",
			model:       fixtureInputModel(),
			getGrafanaConfigsResp: fixtureGrafanaConfigs(func(gc *observability.GrafanaConfigs) {
				gc.GenericOauth = nil
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *observability.ApiUpdateGrafanaConfigsRequest) {
				*request = (*request).UpdateGrafanaConfigsPayload(*fixturePayload(func(payload *observability.UpdateGrafanaConfigsPayload) {
					payload.GenericOauth = nil
				}))
			}),
		},
		{
			description:            "get grafana configs fails",
			model:                  fixtureInputModel(),
			getGrafanaConfigsFails: true,
			isValid:                false,
		},
		{
			description:           "no grafana configs",
			model:                 fixtureInputModel(),
			getGrafanaConfigsResp: nil,
			isValid:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &observabilityClientMocked{
				getGrafanaConfigsFails: tt.getGrafanaConfigsFails,
				getGrafanaConfigsResp:  tt.getGrafanaConfigsResp,
			}
			request, err := buildRequest(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
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
