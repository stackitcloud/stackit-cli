package disable

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &argus.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

type argusClientMocked struct {
	getGrafanaConfigsFails bool
	getGrafanaConfigsResp  *argus.GrafanaConfigs
}

func (c *argusClientMocked) GetInstanceExecute(ctx context.Context, instanceId, projectId string) (*argus.GetInstanceResponse, error) {
	return testClient.GetInstanceExecute(ctx, instanceId, projectId)
}

func (c *argusClientMocked) UpdateGrafanaConfigs(ctx context.Context, instanceId, projectId string) argus.ApiUpdateGrafanaConfigsRequest {
	return testClient.UpdateGrafanaConfigs(ctx, instanceId, projectId)
}

func (c *argusClientMocked) GetGrafanaConfigsExecute(_ context.Context, _, _ string) (*argus.GrafanaConfigs, error) {
	if c.getGrafanaConfigsFails {
		return nil, fmt.Errorf("get payload failed")
	}
	return c.getGrafanaConfigsResp, nil
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:  testProjectId,
		instanceIdFlag: testInstanceId,
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

func fixtureGrafanaConfigs(mods ...func(gc *argus.GrafanaConfigs)) *argus.GrafanaConfigs {
	gc := argus.GrafanaConfigs{
		GenericOauth: &argus.GrafanaOauth{
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

func fixtureRequest(mods ...func(request *argus.ApiUpdateGrafanaConfigsRequest)) argus.ApiUpdateGrafanaConfigsRequest {
	request := testClient.UpdateGrafanaConfigs(testCtx, testInstanceId, testProjectId)
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
			cmd := NewCmd(nil)
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

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description            string
		model                  *inputModel
		getGrafanaConfigsFails bool
		getGrafanaConfigsResp  *argus.GrafanaConfigs
		isValid                bool
		expectedRequest        argus.ApiUpdateGrafanaConfigsRequest
	}{
		{
			description:            "base",
			model:                  fixtureInputModel(),
			getGrafanaConfigsFails: false,
			getGrafanaConfigsResp:  fixtureGrafanaConfigs(),
			isValid:                true,
			expectedRequest: fixtureRequest().UpdateGrafanaConfigsPayload(
				argus.UpdateGrafanaConfigsPayload{
					GenericOauth:     argusUtils.ToPayloadGenericOAuth(fixtureGrafanaConfigs().GenericOauth),
					PublicReadAccess: fixtureGrafanaConfigs().PublicReadAccess,
					UseStackitSso:    utils.Ptr(false),
				}),
		},
		{
			description:            "nil generic oauth",
			model:                  fixtureInputModel(),
			getGrafanaConfigsFails: false,
			getGrafanaConfigsResp: fixtureGrafanaConfigs(func(gc *argus.GrafanaConfigs) {
				gc.GenericOauth = nil
			}),
			isValid: true,
			expectedRequest: fixtureRequest().UpdateGrafanaConfigsPayload(
				argus.UpdateGrafanaConfigsPayload{
					GenericOauth:     nil,
					PublicReadAccess: fixtureGrafanaConfigs().PublicReadAccess,
					UseStackitSso:    utils.Ptr(false),
				}),
		},
		{
			description:            "get grafana configs fails",
			model:                  fixtureInputModel(),
			getGrafanaConfigsFails: true,
			isValid:                false,
		},
		{
			description:            "no grafana configs",
			model:                  fixtureInputModel(),
			getGrafanaConfigsFails: false,
			getGrafanaConfigsResp:  nil,
			isValid:                false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &argusClientMocked{
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
