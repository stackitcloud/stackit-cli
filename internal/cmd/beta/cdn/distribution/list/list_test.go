package list

import (
	"bytes"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

type testCtxKey struct{}

var testProjectId = uuid.NewString()
var testClient = &cdn.APIClient{}
var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")

const (
	testNextPageID = "next-page-id-123"
	testID         = "dist-1"
	testStatus     = cdn.DISTRIBUTIONSTATUS_ACTIVE
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	m := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		SortBy: "createdAt",
	}
	for _, mod := range mods {
		mod(m)
	}
	return m
}

func fixtureRequest(mods ...func(request cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
	request := testClient.ListDistributions(testCtx, testProjectId)
	request = request.PageSize(100)
	request = request.SortBy("createdAt")
	for _, mod := range mods {
		request = mod(request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description string
		argValues   []string
		flagValues  map[string]string
		isValid     bool
		expected    *inputModel
	}{
		{
			description: "base",
			flagValues:  fixtureFlagValues(),
			isValid:     true,
			expected:    fixtureInputModel(),
		},
		{
			description: "no project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "sort by id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "id"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "id"
			}),
		},
		{
			description: "sort by origin-url",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "originUrl"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "originUrl"
			}),
		},
		{
			description: "sort by status",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "status"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "status"
			}),
		},
		{
			description: "sort by created",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "createdAt"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "createdAt"
			}),
		},
		{
			description: "sort by updated",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "updatedAt"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "updatedAt"
			}),
		},
		{
			description: "sort by originUrlRelated",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "originUrlRelated"
			}),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "originUrlRelated"
			}),
		},
		{
			description: "invalid sort by",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sortByFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "missing sort by uses default",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, sortByFlag)
				},
			),
			isValid: true,
			expected: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "createdAt"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expected, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description string
		inputModel  *inputModel
		nextPageID  *string
		expected    cdn.ApiListDistributionsRequest
	}{
		{
			description: "base",
			inputModel:  fixtureInputModel(),
			expected:    fixtureRequest(),
		},
		{
			description: "sort by updatedAt",
			inputModel: fixtureInputModel(func(model *inputModel) {
				model.SortBy = "updatedAt"
			}),
			expected: fixtureRequest(func(req cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
				return req.SortBy("updatedAt")
			}),
		},
		{
			description: "with next page id",
			inputModel:  fixtureInputModel(),
			nextPageID:  utils.Ptr(testNextPageID),
			expected: fixtureRequest(func(req cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
				return req.PageIdentifier(testNextPageID)
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := buildRequest(testCtx, tt.inputModel, testClient, tt.nextPageID, maxPageSize)
			diff := cmp.Diff(req, tt.expected,
				cmp.AllowUnexported(tt.expected),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Errorf("buildRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type testResponse struct {
	statusCode int
	body       cdn.ListDistributionsResponse
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

func fixtureDistributions(count int) []cdn.Distribution {
	distributions := make([]cdn.Distribution, count)
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("dist-%d", i+1)
		distributions[i] = cdn.Distribution{
			Id: &id,
		}
	}
	return distributions
}

func TestFetchDistributions(t *testing.T) {
	tests := []struct {
		description string
		limit       int32
		responses   []testResponse
		expected    []cdn.Distribution
		fails       bool
	}{
		{
			description: "no distributions",
			responses: []testResponse{
				fixtureTestResponse(),
			},
			expected: nil,
		},
		{
			description: "single distribution, single page",
			responses: []testResponse{
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.body.Distributions = &[]cdn.Distribution{
							{Id: utils.Ptr("dist-1")},
						}
					},
				),
			},
			expected: []cdn.Distribution{
				{Id: utils.Ptr("dist-1")},
			},
		},
		{
			description: "multiple distributions, multiple pages",
			responses: []testResponse{
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.body.NextPageIdentifier = utils.Ptr(testNextPageID)
						resp.body.Distributions = &[]cdn.Distribution{
							{Id: utils.Ptr("dist-1")},
						}
					},
				),
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.body.Distributions = &[]cdn.Distribution{
							{Id: utils.Ptr("dist-2")},
						}
					},
				),
			},
			expected: []cdn.Distribution{
				{Id: utils.Ptr("dist-1")},
				{Id: utils.Ptr("dist-2")},
			},
		},
		{
			description: "API error",
			responses: []testResponse{
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.statusCode = 500
					},
				),
			},
			fails: true,
		},
		{
			description: "API error on second page",
			responses: []testResponse{
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.body.NextPageIdentifier = utils.Ptr(testNextPageID)
						resp.body.Distributions = &[]cdn.Distribution{
							{Id: utils.Ptr("dist-1")},
						}
					},
				),
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.statusCode = 500
					},
				),
			},
			fails: true,
		},
		{
			description: "limit across 2 pages",
			limit:       110,
			responses: []testResponse{
				fixtureTestResponse(
					func(resp *testResponse) {
						resp.body.NextPageIdentifier = utils.Ptr(testNextPageID)
						distributions := fixtureDistributions(100)
						resp.body.Distributions = &distributions
					},
				),
				fixtureTestResponse(
					func(resp *testResponse) {
						distributions := fixtureDistributions(10)
						resp.body.Distributions = &distributions
					},
				),
			},
			expected: slices.Concat(fixtureDistributions(100), fixtureDistributions(10)),
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
			client, err := cdn.NewAPIClient(
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
			got, err := fetchDistributions(testCtx, model, client)
			if err != nil {
				if !tt.fails {
					t.Fatalf("fetchDistributions() unexpected error: %v", err)
				}
				return
			}
			if callCount != len(tt.responses) {
				t.Errorf("fetchDistributions() expected %d calls, got %d", len(tt.responses), callCount)
			}
			diff := cmp.Diff(got, tt.expected)
			if diff != "" {
				t.Errorf("fetchDistributions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		description   string
		outputFormat  string
		distributions []cdn.Distribution
		expected      string
	}{
		{
			description:   "no distributions",
			outputFormat:  "json",
			distributions: []cdn.Distribution{},
			expected: `[]
`,
		},
		{
			description:  "no distributions nil slice",
			outputFormat: "json",
			expected: `[]
`,
		},
		{
			description:  "single distribution",
			outputFormat: "table",
			distributions: []cdn.Distribution{
				{
					Id: utils.Ptr(testID),
					Config: &cdn.Config{
						Regions: &[]cdn.Region{
							cdn.REGION_EU,
							cdn.REGION_AF,
						},
					},
					Status: utils.Ptr(testStatus),
				},
			},
			expected: `
 ID     │ REGIONS │ STATUS 
────────┼─────────┼────────
 dist-1 │ EU, AF  │ ACTIVE 

`,
		},
		{
			description:  "no distributions, table format",
			outputFormat: "table",
			expected:     "No CDN distributions found\n",
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			p.Cmd.SetOut(buffer)
			if err := outputResult(p, tt.outputFormat, tt.distributions); err != nil {
				t.Fatalf("outputResult: %v", err)
			}
			if buffer.String() != tt.expected {
				t.Errorf("want:\n%s\ngot:\n%s", tt.expected, buffer.String())
			}
		})
	}
}
