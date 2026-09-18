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

	testDisplayName = "display-name"
	testDescription = "description"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		instanceIdFlag:  testInstanceId,
		displayNameFlag: testDisplayName,
		descriptionFlag: testDescription,
		ttlFlag:         "30",
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

		InstanceId:  testInstanceId,
		DisplayName: testDisplayName,
		Description: new(testDescription),
		Ttl:         new(int32(30)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiCreateAccessTokenRequest)) telemetryrouter.ApiCreateAccessTokenRequest {
	request := testClient.DefaultAPI.CreateAccessToken(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.CreateAccessTokenPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *telemetryrouter.CreateAccessTokenPayload)) telemetryrouter.CreateAccessTokenPayload {
	payload := telemetryrouter.CreateAccessTokenPayload{
		DisplayName: testDisplayName,
		Description: new(testDescription),
		Ttl:         *telemetryrouter.NewNullableInt32(new(int32(30))),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			description: "only required flags",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, ttlFlag)
				delete(flagValues, descriptionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Ttl = nil
				model.Description = nil
			}),
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
			description: "display name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "ttl invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[ttlFlag] = "invalid-integer"
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
		expectedRequest telemetryrouter.ApiCreateAccessTokenRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no optional values",
			model: fixtureInputModel(func(model *inputModel) {
				model.Ttl = nil
				model.Description = nil
			}),
			expectedRequest: fixtureRequest(func(request *telemetryrouter.ApiCreateAccessTokenRequest) {
				*request = request.CreateAccessTokenPayload(telemetryrouter.CreateAccessTokenPayload{
					DisplayName: testDisplayName,
					Description: nil,
					Ttl:         *telemetryrouter.NewNullableInt32(nil),
				})
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(tt.expectedRequest, request,
				cmp.AllowUnexported(tt.expectedRequest, telemetryrouter.NullableInt32{}),
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
		model       *inputModel
		accessToken *telemetryrouter.CreateAccessTokenResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				model: &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				accessToken: new(telemetryrouter.CreateAccessTokenResponse{
					Id:          uuid.NewString(),
					DisplayName: "Token",
					AccessToken: "Secret access token",
					CreatorId:   uuid.NewString(),
					Status:      telemetryrouter.ACCESSTOKENBASERESPONSESTATUS_ACTIVE,
				}),
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{OutputFormat: print.JSONOutputFormat}},
				accessToken: new(telemetryrouter.CreateAccessTokenResponse{}),
			},
			wantErr: false,
		},
		{
			name: "empty response",
			args: args{
				model:       &inputModel{GlobalFlagModel: &globalflags.GlobalFlagModel{}},
				accessToken: nil,
			},
			wantErr: true,
		},
		{
			name: "nil model",
			args: args{
				model:       nil,
				accessToken: new(telemetryrouter.CreateAccessTokenResponse{}),
			},
			wantErr: true,
		},
		{
			name: "nil global flag model",
			args: args{
				model:       &inputModel{GlobalFlagModel: nil},
				accessToken: new(telemetryrouter.CreateAccessTokenResponse{}),
			},
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, "label", tt.args.accessToken); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
