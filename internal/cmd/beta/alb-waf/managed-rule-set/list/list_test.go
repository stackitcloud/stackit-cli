package list

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "test")
	testProjectId  = uuid.NewString()
	testRegion     = "eu01"
	testClient     = &albwaf.APIClient{DefaultAPI: &albwaf.DefaultAPIService{}}
	testNextPageID = "next-page-id-123"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
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
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request albwaf.ApiListManagedRuleSetsRequest) albwaf.ApiListManagedRuleSetsRequest) albwaf.ApiListManagedRuleSetsRequest {
	request := testClient.DefaultAPI.ListManagedRuleSets(testCtx, testProjectId, testRegion)
	request = request.PageSize("100")
	for _, mod := range mods {
		request = mod(request)
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
			description: "invalid limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		nextPageID      string
		pageLimit       int64
		expectedRequest albwaf.ApiListManagedRuleSetsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			pageLimit:       100,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with next page id",
			model:       fixtureInputModel(),
			nextPageID:  testNextPageID,
			pageLimit:   100,
			expectedRequest: fixtureRequest(func(request albwaf.ApiListManagedRuleSetsRequest) albwaf.ApiListManagedRuleSetsRequest {
				return request.PageId(testNextPageID)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, tt.nextPageID, tt.pageLimit)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, albwaf.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("data does not match: %s", diff)
			}
		})
	}
}

type testResponse struct {
	statusCode int
	body       albwaf.ListManagedRuleSetResponse
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

func fixtureItems(count int) []albwaf.GetLimitedManagedRuleSetResponse {
	items := make([]albwaf.GetLimitedManagedRuleSetResponse, count)
	for i := range count {
		items[i] = albwaf.GetLimitedManagedRuleSetResponse{
			Name:                 fmt.Sprintf("managed-rule-set-%d", i+1),
			Type:                 albwaf.TYPE_TYPE_OWASP_CRS,
			Version:              "v1.0.0",
			AdditionalProperties: map[string]interface{}{},
		}
	}
	return items
}

func TestFetchManagedRuleSets(t *testing.T) {
	tests := []struct {
		description string
		limit       int64
		responses   []testResponse
		expected    []albwaf.GetLimitedManagedRuleSetResponse
		fails       bool
	}{
		{
			description: "no items",
			responses: []testResponse{
				fixtureTestResponse(),
			},
			expected: nil,
		},
		{
			description: "single item, single page",
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.Items = fixtureItems(1)
				}),
			},
			expected: fixtureItems(1),
		},
		{
			description: "multiple pages",
			responses: []testResponse{
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.NextPageId = utils.Ptr(testNextPageID)
					resp.body.Items = fixtureItems(1)
				}),
				fixtureTestResponse(func(resp *testResponse) {
					resp.body.Items = fixtureItems(1)
				}),
			},
			expected: slices.Concat(fixtureItems(1), fixtureItems(1)),
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
			client, err := albwaf.NewAPIClient(
				sdkConfig.WithEndpoint(server.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("failed to create test client: %v", err)
			}
			var mods []func(m *inputModel)
			if tt.limit > 0 {
				mods = append(mods, func(m *inputModel) {
					m.Limit = utils.Ptr(tt.limit)
				})
			}
			model := fixtureInputModel(mods...)
			got, err := fetchManagedRuleSets(testCtx, model, client)
			if err != nil {
				if !tt.fails {
					t.Fatalf("fetchManagedRuleSets() unexpected error: %v", err)
				}
				return
			}
			if callCount != len(tt.responses) {
				t.Errorf("fetchManagedRuleSets() expected %d calls, got %d", len(tt.responses), callCount)
			}
			diff := cmp.Diff(got, tt.expected)
			if diff != "" {
				t.Errorf("fetchManagedRuleSets() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description string
		items       []albwaf.GetLimitedManagedRuleSetResponse
	}{
		{
			description: "no items",
			items:       []albwaf.GetLimitedManagedRuleSetResponse{},
		},
		{
			description: "single item",
			items:       fixtureItems(1),
		},
	}

	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if err := outputResult(params.Printer, "", tt.items); err != nil {
				t.Fatalf("outputResult: %v", err)
			}
		})
	}
}
