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

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &telemetryrouter.APIClient{DefaultAPI: &telemetryrouter.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceIdFlag:            testInstanceId,
		limitFlag:                 "10",
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
		InstanceId: testInstanceId,
		Limit:      new(int64(10)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *telemetryrouter.ApiListDestinationsRequest)) telemetryrouter.ApiListDestinationsRequest {
	request := testClient.DefaultAPI.ListDestinations(testCtx, testProjectId, testRegion, testInstanceId)
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
		expectedRequest telemetryrouter.ApiListDestinationsRequest
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
	body       telemetryrouter.ListDestinationsResponse
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

func fixtureDestinations(count int) []telemetryrouter.DestinationResponse {
	destinations := make([]telemetryrouter.DestinationResponse, count)
	for i := 0; i < count; i++ {
		destinations[i] = telemetryrouter.DestinationResponse{
			Id:          fmt.Sprintf("destination-%d", i+1),
			DisplayName: fmt.Sprintf("destination-%d", i+1),
			Config: telemetryrouter.DestinationConfig{
				ConfigType:           telemetryrouter.DESTINATIONCONFIGTYPE_S3,
				AdditionalProperties: map[string]interface{}{},
			},
			CredentialType:       telemetryrouter.DESTINATIONRESPONSECREDENTIALTYPE_ACCESS_KEY,
			Status:               telemetryrouter.DESTINATIONRESPONSESTATUS_ACTIVE,
			AdditionalProperties: map[string]interface{}{},
		}
	}
	return destinations
}

func TestFetchDestinations(t *testing.T) {
	tests := []struct {
		description string
		limit       int64
		responses   []testResponse
		expected    []telemetryrouter.DestinationResponse
		fails       bool
	}{
		{
			description: "no destinations",
			responses: []testResponse{
				fixtureTestResponse(),
			},
			expected: []telemetryrouter.DestinationResponse{},
		},
		{
			description: "destinations returned",
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.Destinations = fixtureDestinations(3)
				}),
			},
			expected: fixtureDestinations(3),
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
					resp.body.Destinations = fixtureDestinations(5)
				}),
			},
			expected: fixtureDestinations(2),
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

			got, err := fetchDestinations(testCtx, model, client)
			if err != nil {
				if !tt.fails {
					t.Fatalf("fetchDestinations() unexpected error: %v", err)
				}
				return
			}
			if tt.fails {
				t.Fatalf("fetchDestinations() expected error, got none")
			}
			if callCount != len(tt.responses) {
				t.Errorf("fetchDestinations() expected %d calls, got %d", len(tt.responses), callCount)
			}
			diff := cmp.Diff(got, tt.expected)
			if diff != "" {
				t.Errorf("fetchDestinations() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		instanceLabel string
		destinations  []telemetryrouter.DestinationResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "empty destinations slice",
			args: args{
				destinations: []telemetryrouter.DestinationResponse{},
			},
			wantErr: false,
		},
		{
			name: "empty destination in destinations slice",
			args: args{
				destinations: []telemetryrouter.DestinationResponse{{}},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.destinations, tt.args.instanceLabel); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
