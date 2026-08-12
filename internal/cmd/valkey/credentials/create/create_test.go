package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &valkey.APIClient{DefaultAPI: &valkey.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
		instanceIdFlag:         testInstanceId,
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
			Region:    testRegion,
		},
		InstanceId:   testInstanceId,
		ShowPassword: false,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *valkey.ApiCreateCredentialsRequest)) valkey.ApiCreateCredentialsRequest {
	request := testClient.DefaultAPI.CreateCredentials(testCtx, testProjectId, testRegion, testInstanceId)
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
			description: "with show-password",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[showPasswordFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ShowPassword = true
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
			}),
			isValid: false,
		},
		{
			description: "instance id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, []string{}, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest valkey.ApiCreateCredentialsRequest
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
				cmp.AllowUnexported(tt.expectedRequest, valkey.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name          string
		model         *inputModel
		instanceLabel string
		resp          *valkey.CredentialsResponse
		wantErr       bool
	}{
		{
			name:          "nil response",
			model:         fixtureInputModel(),
			instanceLabel: "example-instance",
			resp:          nil,
			wantErr:       true,
		},
		{
			name:          "response without raw credentials",
			model:         fixtureInputModel(),
			instanceLabel: "example-instance",
			resp: &valkey.CredentialsResponse{
				Id:  "creds-id",
				Uri: "redis://host:6379",
			},
			wantErr: false,
		},
		{
			name:          "response with raw credentials, password hidden",
			model:         fixtureInputModel(),
			instanceLabel: "example-instance",
			resp: &valkey.CredentialsResponse{
				Id:  "creds-id",
				Uri: "redis://host:6379",
				Raw: &valkey.RawCredentials{
					Credentials: valkey.Credentials{
						Host:     "host",
						Username: "user",
						Password: "secret",
						Port:     new(int32(6379)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "response with raw credentials, password shown",
			model: fixtureInputModel(func(model *inputModel) {
				model.ShowPassword = true
			}),
			instanceLabel: "example-instance",
			resp: &valkey.CredentialsResponse{
				Id:  "creds-id",
				Uri: "redis://host:6379",
				Raw: &valkey.RawCredentials{
					Credentials: valkey.Credentials{
						Host:     "host",
						Username: "user",
						Password: "secret",
						Port:     new(int32(6379)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "async mode",
			model: fixtureInputModel(func(model *inputModel) {
				model.Async = true
			}),
			instanceLabel: "example-instance",
			resp: &valkey.CredentialsResponse{
				Id:  "creds-id",
				Uri: "redis://host:6379",
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.model, tt.instanceLabel, tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
