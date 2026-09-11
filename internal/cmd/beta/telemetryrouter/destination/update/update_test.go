package update

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx           = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient        = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
	testProjectId     = uuid.NewString()
	testInstanceId    = uuid.NewString()
	testDestinationId = uuid.NewString()
)

var testPayload = telemetryrouter.PartialUpdateDestinationPayload{
	DisplayName: new("my-destination"),
	Description: new("description"),
	Config: &telemetryrouter.DestinationConfigPatch{
		ConfigType: new(telemetryrouter.DESTINATIONCONFIGTYPE_S3),
		S3: &telemetryrouter.DestinationConfigPatchS3{
			Bucket:   new("my-bucket"),
			Endpoint: new("https://object.storage.eu01.onstackit.cloud"),
			AccessKey: &telemetryrouter.DestinationConfigPatchS3AccessKey{ // nolint:gosec // false positive
				Id:     new("access-key-id"),
				Secret: new("access-key-secret"),
			},
		},
	},
}

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testDestinationId,
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
		instanceIdFlag:            testInstanceId,
		payloadFlag: `{
			"displayName": "my-destination",
			"description": "description",
			"config": {
				"configType": "S3",
				"s3": {
					"bucket": "my-bucket",
					"endpoint": "https://object.storage.eu01.onstackit.cloud",
					"accessKey": {
						"id": "access-key-id",
						"secret": "access-key-secret"
					}
				}
			}
		}`,
	}
	for _, mod := range mods {
		mod(flagValues)
	}

	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			Verbosity: globalflags.VerbosityDefault,
			ProjectId: testProjectId,
			Region:    testRegion,
		},
		InstanceId:    testInstanceId,
		DestinationId: testDestinationId,
		Payload:       testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}

	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiPartialUpdateDestinationRequest)) telemetryrouter.ApiPartialUpdateDestinationRequest {
	request := testClient.DefaultAPI.PartialUpdateDestination(testCtx, testProjectId, testRegion, testInstanceId, testDestinationId)

	request = request.PartialUpdateDestinationPayload(testPayload)
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
			description: "instance id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "destination id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "destination id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid json",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
			}),
			isValid: false,
		},
		{
			description: "empty payload",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "{}"
			}),
			isValid: false,
		},
		{
			description: "display name only",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = `{"displayName": "my-destination"}`
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Payload = telemetryrouter.PartialUpdateDestinationPayload{
					DisplayName: new("my-destination"),
				}
			}),
		},
		{
			description: "unknown fields only",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = `{"dispayName": "my-destination"}`
			}),
			isValid: false,
		},
		{
			description: "payload missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, payloadFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithOptions(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, nil, tt.isValid, []testutils.TestingOption{
				testutils.WithCmpOptions(cmp.FilterPath(func(p cmp.Path) bool {
					last := p.Last().String()
					return last == ".AdditionalProperties"
				}, cmp.Ignore())),
			})
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest telemetryrouter.ApiPartialUpdateDestinationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx, telemetryrouter.DefaultAPIService{}),
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
		resp        *telemetryrouter.DestinationResponse
		wantErr     bool
	}{
		{
			description: "nil response",
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			resp:        nil,
			wantErr:     true,
		},
		{
			description: "model is nil",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       nil,
			wantErr:     true,
		},
		{
			description: "global flag nil",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       &inputModel{GlobalFlagModel: nil},
			wantErr:     true,
		},
		{
			description: "default output",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
			wantErr:     false,
		},
		{
			description: "async output",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{Async: true}},
			wantErr:     false,
		},
		{
			description: "json output",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
			wantErr:     false,
		},
		{
			description: "yaml output",
			resp:        &telemetryrouter.DestinationResponse{},
			model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.YAMLOutputFormat}},
			wantErr:     false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := outputResult(params.Printer, tt.model, "label", tt.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
