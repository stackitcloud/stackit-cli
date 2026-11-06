package create

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
	telemetrylink "github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	testRegion      = "eu01"
	testDisplayName = "my-telemetrylink-instance"
	testDescription = "my link description"
	testAccessToken = "test-access-token"
)

type testCtxKey struct{}

var (
	testCtx               = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient            = &telemetrylink.APIClient{DefaultAPI: &telemetrylink.DefaultAPIService{}}
	testResourceId        = uuid.NewString()
	testProjectId         = uuid.NewString()
	testTelemetryRouterId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           testDisplayName,
		descriptionFlag:           testDescription,
		enabledFlag:               "true",
		resourceIdFlag:            testResourceId,
		resourceTypeFlag:          "project",
		telemetryRouterIdFlag:     testTelemetryRouterId,
		accessTokenFlag:           testAccessToken,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// Input Model
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		DisplayName:       testDisplayName,
		Description:       new(testDescription),
		Enabled:           new(true),
		ResourceId:        testResourceId,
		ResourceType:      "project",
		TelemetryRouterId: testTelemetryRouterId,
		AccessToken:       testAccessToken,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureProjectRequest(mods ...func(request *telemetrylink.ApiCreateOrUpdateProjectTelemetryLinkRequest)) telemetrylink.ApiCreateOrUpdateProjectTelemetryLinkRequest {
	request := testClient.DefaultAPI.CreateOrUpdateProjectTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.CreateOrUpdateProjectTelemetryLinkPayload(telemetrylink.CreateOrUpdateProjectTelemetryLinkPayload{
		DisplayName:       testDisplayName,
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: testTelemetryRouterId,
		AccessToken:       testAccessToken,
	})

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureOrganizationRequest(mods ...func(request *telemetrylink.ApiCreateOrUpdateOrganizationTelemetryLinkRequest)) telemetrylink.ApiCreateOrUpdateOrganizationTelemetryLinkRequest {
	request := testClient.DefaultAPI.CreateOrUpdateOrganizationTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.CreateOrUpdateOrganizationTelemetryLinkPayload(telemetrylink.CreateOrUpdateOrganizationTelemetryLinkPayload{
		DisplayName:       testDisplayName,
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: testTelemetryRouterId,
		AccessToken:       testAccessToken,
	})

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureFolderRequest(mods ...func(request *telemetrylink.ApiCreateOrUpdateFolderTelemetryLinkRequest)) telemetrylink.ApiCreateOrUpdateFolderTelemetryLinkRequest {
	request := testClient.DefaultAPI.CreateOrUpdateFolderTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.CreateOrUpdateFolderTelemetryLinkPayload(telemetrylink.CreateOrUpdateFolderTelemetryLinkPayload{
		DisplayName:       testDisplayName,
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: testTelemetryRouterId,
		AccessToken:       testAccessToken,
	})

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
			description: "optional flags omitted",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
		},
		{
			description: "no values provided",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "telemetry router id missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, telemetryRouterIdFlag)
			}),
			isValid: false,
		},
		{
			description: "access token missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, accessTokenFlag)
			}),
			isValid: false,
		},
		{
			description: "resource type missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, resourceTypeFlag)
			}),
			isValid: false,
		},
		{
			description: "resource id missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, resourceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "resource id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[resourceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "resource id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[resourceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "display name missing (required)",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
				return parseInput(p, cmd)
			}, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildProjectRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetrylink.ApiCreateOrUpdateProjectTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureProjectRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			expectedRequest: fixtureProjectRequest().CreateOrUpdateProjectTelemetryLinkPayload(telemetrylink.CreateOrUpdateProjectTelemetryLinkPayload{
				DisplayName:       testDisplayName,
				Description:       nil,
				Enabled:           new(true),
				TelemetryRouterId: testTelemetryRouterId,
				AccessToken:       testAccessToken,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, _ := buildProjectRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, telemetrylink.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildFolderRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetrylink.ApiCreateOrUpdateFolderTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureFolderRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			expectedRequest: fixtureFolderRequest().CreateOrUpdateFolderTelemetryLinkPayload(telemetrylink.CreateOrUpdateFolderTelemetryLinkPayload{
				DisplayName:       testDisplayName,
				Description:       nil,
				Enabled:           new(true),
				TelemetryRouterId: testTelemetryRouterId,
				AccessToken:       testAccessToken,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, _ := buildFolderRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, telemetrylink.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildOrganizationRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetrylink.ApiCreateOrUpdateOrganizationTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureOrganizationRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
			}),
			expectedRequest: fixtureOrganizationRequest().CreateOrUpdateOrganizationTelemetryLinkPayload(telemetrylink.CreateOrUpdateOrganizationTelemetryLinkPayload{
				DisplayName:       testDisplayName,
				Description:       nil,
				Enabled:           new(true),
				TelemetryRouterId: testTelemetryRouterId,
				AccessToken:       testAccessToken,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, _ := buildOrganizationRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, telemetrylink.DefaultAPIService{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		instance    *telemetrylink.TelemetryLinkResponse
		wantErr     bool
	}{
		{
			description: "nil response",
			instance:    nil,
			wantErr:     true,
		},
		{
			description: "model is nil",
			instance:    &telemetrylink.TelemetryLinkResponse{},
			model:       nil,
			wantErr:     true,
		},
		{
			description: "global flag nil",
			instance:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: nil},
			wantErr:     true,
		},
		{
			description: "default output",
			instance:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
		},
		{
			description: "json output",
			instance:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wantErr:     false,
		},
		{
			description: "yaml output",
			instance:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wantErr:     false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(params.Printer, tt.model, true, "label", tt.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
