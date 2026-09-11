package generatepayload

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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

const (
	testFilePath = "example-file"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceIdFlag:            testInstanceId,
		destinationIdFlag:         testDestinationId,
		configTypeFlag:            configTypeOpenTelemetry,
		filePathFlag:              testFilePath,
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
		InstanceId:    testInstanceId,
		DestinationId: new(testDestinationId),
		ConfigType:    configTypeOpenTelemetry,
		FilePath:      new(testFilePath),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiGetDestinationRequest)) telemetryrouter.ApiGetDestinationRequest {
	request := testClient.DefaultAPI.GetDestination(testCtx, testProjectId, testRegion, testInstanceId, testDestinationId)
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
			isValid:     true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{Verbosity: globalflags.VerbosityDefault},
				ConfigType:      configTypeOpenTelemetry,
			},
		},
		{
			description: "file path missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, filePathFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.FilePath = nil
			}),
		},
		{
			description: "destination id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, destinationIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DestinationId = nil
			}),
		},
		{
			description: "config type s3",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[configTypeFlag] = configTypeS3
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ConfigType = configTypeS3
			}),
		},
		{
			description: "config type invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[configTypeFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "instance id missing, destination id provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id missing, destination id provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id and instance id missing, destination id provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
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
		expectedRequest telemetryrouter.ApiGetDestinationRequest
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

func TestOutputCreateResult(t *testing.T) {
	type args struct {
		filePath *string
		payload  *telemetryrouter.CreateDestinationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "empty payload",
			args: args{
				payload: &telemetryrouter.CreateDestinationPayload{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputCreateResult(params.Printer, tt.args.filePath, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("outputCreateResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOutputUpdateResult(t *testing.T) {
	type args struct {
		filePath *string
		payload  *telemetryrouter.UpdateDestinationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "empty payload",
			args: args{
				payload: &telemetryrouter.UpdateDestinationPayload{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputUpdateResult(params.Printer, tt.args.filePath, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("outputUpdateResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
