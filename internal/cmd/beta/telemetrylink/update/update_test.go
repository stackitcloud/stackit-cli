package update

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
		DisplayName:       new(testDisplayName),
		Description:       new(testDescription),
		Enabled:           new(true),
		ResourceId:        testResourceId,
		ResourceType:      "project",
		TelemetryRouterId: new(testTelemetryRouterId),
		AccessToken:       new(testAccessToken),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureProjectRequest(mods ...func(request *telemetrylink.ApiPartialUpdateProjectTelemetryLinkRequest)) telemetrylink.ApiPartialUpdateProjectTelemetryLinkRequest {
	request := testClient.DefaultAPI.PartialUpdateProjectTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.PartialUpdateProjectTelemetryLinkPayload(telemetrylink.PartialUpdateProjectTelemetryLinkPayload{
		DisplayName:       new(testDisplayName),
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: new(testTelemetryRouterId),
		AccessToken:       new(testAccessToken),
	})

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureOrganizationRequest(mods ...func(request *telemetrylink.ApiPartialUpdateOrganizationTelemetryLinkRequest)) telemetrylink.ApiPartialUpdateOrganizationTelemetryLinkRequest {
	request := testClient.DefaultAPI.PartialUpdateOrganizationTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.PartialUpdateOrganizationTelemetryLinkPayload(telemetrylink.PartialUpdateOrganizationTelemetryLinkPayload{
		DisplayName:       new(testDisplayName),
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: new(testTelemetryRouterId),
		AccessToken:       new(testAccessToken),
	})

	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureFolderRequest(mods ...func(request *telemetrylink.ApiPartialUpdateFolderTelemetryLinkRequest)) telemetrylink.ApiPartialUpdateFolderTelemetryLinkRequest {
	request := testClient.DefaultAPI.PartialUpdateFolderTelemetryLink(testCtx, testResourceId, testRegion)
	request = request.PartialUpdateFolderTelemetryLinkPayload(telemetrylink.PartialUpdateFolderTelemetryLinkPayload{
		DisplayName:       new(testDisplayName),
		Description:       new(testDescription),
		Enabled:           new(true),
		TelemetryRouterId: new(testTelemetryRouterId),
		AccessToken:       new(testAccessToken),
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
				delete(flagValues, displayNameFlag)
				delete(flagValues, telemetryRouterIdFlag)
				delete(flagValues, accessTokenFlag)
				delete(flagValues, enabledFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Description = nil
				model.TelemetryRouterId = nil
				model.AccessToken = nil
				model.Enabled = nil
				model.DisplayName = nil
			}),
		},
		{
			description: "no values provided",
			flagValues:  map[string]string{},
			isValid:     false,
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
		expectedRequest telemetrylink.ApiPartialUpdateProjectTelemetryLinkRequest
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
				model.DisplayName = nil
				model.Enabled = nil
				model.TelemetryRouterId = nil
				model.AccessToken = nil
			}),
			expectedRequest: fixtureProjectRequest().PartialUpdateProjectTelemetryLinkPayload(telemetrylink.PartialUpdateProjectTelemetryLinkPayload{
				DisplayName:       nil,
				Description:       nil,
				Enabled:           nil,
				TelemetryRouterId: nil,
				AccessToken:       nil,
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
		expectedRequest telemetrylink.ApiPartialUpdateFolderTelemetryLinkRequest
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
				model.DisplayName = nil
				model.Enabled = nil
				model.TelemetryRouterId = nil
				model.AccessToken = nil
			}),
			expectedRequest: fixtureFolderRequest().PartialUpdateFolderTelemetryLinkPayload(telemetrylink.PartialUpdateFolderTelemetryLinkPayload{
				DisplayName:       nil,
				Description:       nil,
				Enabled:           nil,
				TelemetryRouterId: nil,
				AccessToken:       nil,
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
		expectedRequest telemetrylink.ApiPartialUpdateOrganizationTelemetryLinkRequest
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
				model.DisplayName = nil
				model.Enabled = nil
				model.TelemetryRouterId = nil
				model.AccessToken = nil
			}),
			expectedRequest: fixtureOrganizationRequest().PartialUpdateOrganizationTelemetryLinkPayload(telemetrylink.PartialUpdateOrganizationTelemetryLinkPayload{
				DisplayName:       nil,
				Description:       nil,
				Enabled:           nil,
				TelemetryRouterId: nil,
				AccessToken:       nil,
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
