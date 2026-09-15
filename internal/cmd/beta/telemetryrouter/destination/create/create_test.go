package create

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
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

var testPayload = &telemetryrouter.CreateDestinationPayload{
	DisplayName: "my-destination",
	Description: nil,
	Config: telemetryrouter.DestinationConfig{
		ConfigType: telemetryrouter.DESTINATIONCONFIGTYPE_OPEN_TELEMETRY,
		OpenTelemetry: &telemetryrouter.DestinationConfigOpenTelemetry{
			Uri: "https://otel-collector.example.com:4317",
			BasicAuth: &telemetryrouter.DestinationConfigOpenTelemetryBasicAuth{
				Username: "username",
				Password: "password",
			},
		},
	},
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceIdFlag:            testInstanceId,
		payloadFlag: `{
			"displayName": "my-destination",
			"config": {
				"configType": "OpenTelemetry",
				"openTelemetry": {
					"uri": "https://otel-collector.example.com:4317",
					"basicAuth": {
						"username": "username",
						"password": "password"
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
		InstanceId: testInstanceId,
		Payload:    testPayload,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiCreateDestinationRequest)) telemetryrouter.ApiCreateDestinationRequest {
	request := testClient.DefaultAPI.CreateDestination(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.CreateDestinationPayload(*testPayload)
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
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "default config",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, payloadFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Payload = nil
			}),
		},
		{
			description: "invalid json",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[payloadFlag] = "not json"
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
		expectedRequest telemetryrouter.ApiCreateDestinationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

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
	type args struct {
		model         *inputModel
		instanceLabel string
		destination   *telemetryrouter.DestinationResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				model:         &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				instanceLabel: "label",
				destination: &telemetryrouter.DestinationResponse{
					Id:          uuid.NewString(),
					DisplayName: "my-destination",
					Status:      telemetryrouter.DESTINATIONRESPONSESTATUS_ACTIVE,
				},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model:         &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
				instanceLabel: "label",
				destination:   &telemetryrouter.DestinationResponse{},
			},
			wantErr: false,
		},
		{
			name: "empty response",
			args: args{
				model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				destination: nil,
			},
			wantErr: true,
		},
		{
			name: "nil model",
			args: args{
				model:       nil,
				destination: &telemetryrouter.DestinationResponse{},
			},
			wantErr: true,
		},
		{
			name: "nil global flag model",
			args: args{
				model:       &inputModel{GlobalFlagModel: nil},
				destination: &telemetryrouter.DestinationResponse{},
			},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.instanceLabel, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
