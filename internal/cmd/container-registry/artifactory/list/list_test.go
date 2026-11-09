package list

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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

func fixtureRequest(mods ...func(request *registry.ApiListArtifactoriesRequest)) registry.ApiListArtifactoriesRequest {
	request := testClient.DefaultAPI.ListArtifactories(testCtx, testProjectId, testRegion).PageSize(maxPageSize)
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
		pageToken       string
		pageSize        int32
		expectedRequest registry.ApiListArtifactoriesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			pageSize:        maxPageSize,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with page token",
			model:       fixtureInputModel(),
			pageToken:   "next-token",
			pageSize:    maxPageSize,
			expectedRequest: fixtureRequest(func(request *registry.ApiListArtifactoriesRequest) {
				*request = request.NextPageToken("next-token")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, tt.pageToken, tt.pageSize)

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
		projectLabel  string
		artifactories []registry.Artifactory
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty artifactories",
			args: args{
				projectLabel:  "example-project",
				artifactories: []registry.Artifactory{},
			},
			wantErr: false,
		},
		{
			name: "single item",
			args: args{
				projectLabel: "example-project",
				artifactories: []registry.Artifactory{
					{
						Id:    uuid.NewString(),
						Name:  "my-registry",
						State: registry.ARTIFACTORYSTATE_ACTIVE,
						Url:   "https://example.com",
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, "", tt.args.projectLabel, tt.args.artifactories); (err != nil) != tt.wantErr {
				t.Errorf("TestOutputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchArtifactories(t *testing.T) {
	tests := []struct {
		description   string
		model         *inputModel
		mockResponses []*registry.ArtifactoryList
		mockErrors    []error
		expectedItems int
		wantErr       bool
	}{
		{
			description: "single page, no limit",
			model:       fixtureInputModel(),
			mockResponses: []*registry.ArtifactoryList{
				{
					Artifactories: []registry.Artifactory{
						{Id: "1", Name: "reg-1"},
						{Id: "2", Name: "reg-2"},
					},
				},
			},
			expectedItems: 2,
			wantErr:       false,
		},
		{
			description: "multi page with pagination",
			model:       fixtureInputModel(),
			mockResponses: []*registry.ArtifactoryList{
				{
					Artifactories: []registry.Artifactory{
						{Id: "1", Name: "reg-1"},
					},
					NextPageToken: utils.Ptr("token-1"),
				},
				{
					Artifactories: []registry.Artifactory{
						{Id: "2", Name: "reg-2"},
					},
				},
			},
			expectedItems: 2,
			wantErr:       false,
		},
		{
			description: "with limit smaller than available",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(1))
			}),
			mockResponses: []*registry.ArtifactoryList{
				{
					Artifactories: []registry.Artifactory{
						{Id: "1", Name: "reg-1"},
						{Id: "2", Name: "reg-2"},
					},
				},
			},
			expectedItems: 1,
			wantErr:       false,
		},
		{
			description: "api call fails",
			model:       fixtureInputModel(),
			mockErrors: []error{
				fmt.Errorf("server error"),
			},
			wantErr: true,
		},
		{
			description: "nil response",
			model:       fixtureInputModel(),
			mockResponses: []*registry.ArtifactoryList{
				nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			callIndex := 0
			mock := registry.DefaultAPIServiceMock{
				ListArtifactoriesExecuteMock: utils.Ptr(func(_ registry.ApiListArtifactoriesRequest) (*registry.ArtifactoryList, error) {
					idx := callIndex
					callIndex++
					if idx < len(tt.mockErrors) && tt.mockErrors[idx] != nil {
						return nil, tt.mockErrors[idx]
					}
					if idx < len(tt.mockResponses) {
						return tt.mockResponses[idx], nil
					}
					return &registry.ArtifactoryList{}, nil
				}),
			}
			client := &registry.APIClient{DefaultAPI: mock}

			got, err := fetchArtifactories(testCtx, tt.model, client)
			if (err != nil) != tt.wantErr {
				t.Fatalf("fetchArtifactories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.expectedItems {
				t.Errorf("fetchArtifactories() returned %d items, want %d", len(got), tt.expectedItems)
			}
		})
	}
}
