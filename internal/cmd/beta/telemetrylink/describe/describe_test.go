package describe

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
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &telemetrylink.APIClient{DefaultAPI: &telemetrylink.DefaultAPIService{}}
	testResourceId = uuid.NewString()
)

// Flags
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.RegionFlag: testRegion,
		resourceIdFlag:         testResourceId,
		resourceTypeFlag:       "project",
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		ResourceId:   testResourceId,
		ResourceType: "project",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// Request
func fixtureProjectRequest(mods ...func(request *telemetrylink.ApiGetProjectTelemetryLinkRequest)) telemetrylink.ApiGetProjectTelemetryLinkRequest {
	request := testClient.DefaultAPI.GetProjectTelemetryLink(testCtx, testResourceId, testRegion)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureOrganizationRequest(mods ...func(request *telemetrylink.ApiGetOrganizationTelemetryLinkRequest)) telemetrylink.ApiGetOrganizationTelemetryLinkRequest {
	request := testClient.DefaultAPI.GetOrganizationTelemetryLink(testCtx, testResourceId, testRegion)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureFolderRequest(mods ...func(request *telemetrylink.ApiGetFolderTelemetryLinkRequest)) telemetrylink.ApiGetFolderTelemetryLinkRequest {
	request := testClient.DefaultAPI.GetFolderTelemetryLink(testCtx, testResourceId, testRegion)
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
		expectedRequest telemetrylink.ApiGetProjectTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureProjectRequest(),
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
		expectedRequest telemetrylink.ApiGetFolderTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureFolderRequest(),
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
		expectedRequest telemetrylink.ApiGetOrganizationTelemetryLinkRequest
	}{
		{
			description:     "base case",
			model:           fixtureInputModel(),
			expectedRequest: fixtureOrganizationRequest(),
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
		response    *telemetrylink.TelemetryLinkResponse
		wantErr     bool
	}{
		{
			description: "nil response",
			response:    nil,
			wantErr:     true,
		},
		{
			description: "model is nil",
			response:    &telemetrylink.TelemetryLinkResponse{},
			model:       nil,
			wantErr:     true,
		},
		{
			description: "global flag nil",
			response:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: nil},
			wantErr:     true,
		},
		{
			description: "default output",
			response:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
		},
		{
			description: "json output",
			response:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wantErr:     false,
		},
		{
			description: "yaml output",
			response:    &telemetrylink.TelemetryLinkResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wantErr:     false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(params.Printer, tt.model, tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
