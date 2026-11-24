package list

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &registry.APIClient{DefaultAPI: &registry.DefaultAPIService{}}
	testProjectId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
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
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *registry.ApiListReplicationAdaptersRequest)) registry.ApiListReplicationAdaptersRequest {
	request := testClient.DefaultAPI.ListReplicationAdapters(testCtx, testProjectId, testRegion)
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
			description: "with limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "10"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
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
		{
			description: "limit invalid: zero",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "limit invalid: negative",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "-1"
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
		expectedRequest registry.ApiListReplicationAdaptersRequest
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
				cmp.AllowUnexported(tt.expectedRequest, registry.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		adapters     []registry.ReplicationAdapter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty adapters",
			args: args{
				projectLabel: "my-project",
				adapters:     []registry.ReplicationAdapter{},
			},
			wantErr: false,
		},
		{
			name: "single item",
			args: args{
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple items",
			args: args{
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
					{
						RegistryType:   "harbor",
						DisplayName:    "Harbor",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "json output format",
			args: args{
				outputFormat: print.JSONOutputFormat,
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.adapters)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchReplicationAdapters(t *testing.T) {
	tests := []struct {
		description   string
		model         *inputModel
		nextPageToken *string
		statusCode    int
		expectedLen   int
		isValid       bool
	}{
		{
			description: "single page",
			model:       fixtureInputModel(),
			statusCode:  http.StatusOK,
			expectedLen: 3,
			isValid:     true,
		},
		{
			description:   "next page token is ignored",
			model:         fixtureInputModel(),
			nextPageToken: utils.Ptr("next-page"),
			statusCode:    http.StatusOK,
			expectedLen:   3,
			isValid:       true,
		},
		{
			description: "with limit",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(2))
			}),
			nextPageToken: utils.Ptr("next-page"),
			statusCode:    http.StatusOK,
			expectedLen:   2,
			isValid:       true,
		},
		{
			description: "api error",
			model:       fixtureInputModel(),
			statusCode:  http.StatusInternalServerError,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var calls atomic.Int32
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.statusCode != http.StatusOK {
					return
				}

				resp := registry.ReplicationAdapterList{
					ReplicationAdapters: []registry.ReplicationAdapter{
						{RegistryType: "docker-hub", DisplayName: "Docker Hub"},
						{RegistryType: "harbor", DisplayName: "Harbor"},
						{RegistryType: "quay", DisplayName: "Quay"},
					},
					NextPageToken: tt.nextPageToken,
				}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					t.Errorf("encode response: %v", err)
				}
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			client, err := registry.NewAPIClient(
				sdkConfig.WithEndpoint(server.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("init client: %v", err)
			}

			adapters, err := fetchReplicationAdapters(testCtx, tt.model, client)
			if !tt.isValid {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(adapters) != tt.expectedLen {
				t.Errorf("expected %d adapters, got %d", tt.expectedLen, len(adapters))
			}
			if got := calls.Load(); got != 1 {
				t.Errorf("expected exactly 1 API call, got %d", got)
			}
		})
	}
}
