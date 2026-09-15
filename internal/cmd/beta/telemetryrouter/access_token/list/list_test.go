package list

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
)

const (
	testRegion = "eu01"
)

type testCtxKey struct{}

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}

	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,

		instanceIdFlag: testInstanceId,
		limitFlag:      "10",
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
		Limit:      new(int64(10)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiListAccessTokensRequest)) telemetryrouter.ApiListAccessTokensRequest {
	request := testClient.DefaultAPI.ListAccessTokens(testCtx, testProjectId, testRegion, testInstanceId)
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
			description: "no flag values",
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
			description: "limit invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
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
		expectedRequest telemetryrouter.ApiListAccessTokensRequest
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

type testResponse struct {
	statusCode int
	body       telemetryrouter.ListAccessTokensResponse
}

func fixtureTestResponse(mods ...func(resp *testResponse)) testResponse {
	resp := testResponse{
		statusCode: 200,
	}
	for _, mod := range mods {
		mod(&resp)
	}
	return resp
}

func fixtureAccessTokens(count int) []telemetryrouter.GetAccessTokenResponse {
	tokens := make([]telemetryrouter.GetAccessTokenResponse, count)
	for i := 0; i < count; i++ {
		tokens[i] = telemetryrouter.GetAccessTokenResponse{
			Id:                   fmt.Sprintf("token-%d", i+1),
			DisplayName:          fmt.Sprintf("token-%d", i+1),
			CreatorId:            testProjectId,
			Status:               telemetryrouter.ACCESSTOKENBASERESPONSESTATUS_ACTIVE,
			ExpirationTime:       *telemetryrouter.NewNullableTime(nil),
			AdditionalProperties: map[string]interface{}{},
		}
	}
	return tokens
}

func TestFetchAccessTokens(t *testing.T) {
	tests := []struct {
		description string
		limit       int64
		responses   []testResponse
		expected    []telemetryrouter.GetAccessTokenResponse
		fails       bool
	}{
		{
			description: "no access tokens",
			responses: []testResponse{
				fixtureTestResponse(),
			},
			expected: []telemetryrouter.GetAccessTokenResponse{},
		},
		{
			description: "access tokens returned",
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.AccessTokens = fixtureAccessTokens(3)
				}),
			},
			expected: fixtureAccessTokens(3),
		},
		{
			description: "API error",
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.statusCode = 500
				}),
			},
			fails: true,
		},
		{
			description: "limit truncates within a single page",
			limit:       2,
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.AccessTokens = fixtureAccessTokens(5)
				}),
			},
			expected: fixtureAccessTokens(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			callCount := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resp := tt.responses[callCount]
				callCount++
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(resp.statusCode)
				bs, err := json.Marshal(resp.body)
				if err != nil {
					t.Fatalf("marshal: %v", err)
				}
				_, err = w.Write(bs)
				if err != nil {
					t.Fatalf("write: %v", err)
				}
			})
			server := httptest.NewServer(handler)
			defer server.Close()
			client, err := telemetryrouter.NewAPIClient(
				sdkConfig.WithEndpoint(server.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("failed to create test client: %v", err)
			}
			model := fixtureInputModel(func(m *inputModel) {
				if tt.limit > 0 {
					m.Limit = utils.Ptr(tt.limit)
				} else {
					m.Limit = nil
				}
			})

			got, err := fetchAccessTokens(testCtx, model, client)
			if err != nil {
				if !tt.fails {
					t.Fatalf("fetchAccessTokens() unexpected error: %v", err)
				}
				return
			}
			if tt.fails {
				t.Fatalf("fetchAccessTokens() expected error, got none")
			}
			if callCount != len(tt.responses) {
				t.Errorf("fetchAccessTokens() expected %d calls, got %d", len(tt.responses), callCount)
			}
			diff := cmp.Diff(got, tt.expected, cmpopts.EquateComparable(telemetryrouter.NullableTime{}))
			if diff != "" {
				t.Errorf("fetchAccessTokens() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		accessTokens  []telemetryrouter.GetAccessTokenResponse
		instanceLabel string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				accessTokens: []telemetryrouter.GetAccessTokenResponse{
					{
						Id:          uuid.NewString(),
						DisplayName: "Token",
						CreatorId:   uuid.NewString(),
						Status:      telemetryrouter.ACCESSTOKENBASERESPONSESTATUS_ACTIVE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "set empty access token",
			args: args{
				accessTokens: []telemetryrouter.GetAccessTokenResponse{
					{},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.accessTokens, tt.args.instanceLabel); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
