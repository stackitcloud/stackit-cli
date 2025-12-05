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
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
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

func flagSortBy(sortBy string) func(m map[string]string) {
	return func(m map[string]string) {
		m[sortByFlag] = sortBy
	}
}

func flagProjectId(id *string) func(m map[string]string) {
	return func(m map[string]string) {
		if id == nil {
			delete(m, globalflags.ProjectIdFlag)
		} else {
			m[globalflags.ProjectIdFlag] = *id
		}
	}
}

func fixtureInputModel(mods ...func(m *inputModel)) *inputModel {
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

func inputSortBy(sortBy string) func(m *inputModel) {
	return func(m *inputModel) {
		m.SortBy = sortBy
	}
}

func fixtureRequest(mods ...func(r cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
	r := testClient.ListDistributions(testCtx, testProjectId)
	r = r.PageSize(100)
	r = r.SortBy("createdAt")
	for _, mod := range mods {
		r = mod(r)
	}
	return r
}

func requestSortBy(sortBy string) func(r cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
	return func(r cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
		return r.SortBy(sortBy)
	}
}

func requestNextPageID(nextPageID string) func(r cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
	return func(r cdn.ApiListDistributionsRequest) cdn.ApiListDistributionsRequest {
		return r.PageIdentifier(nextPageID)
	}
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
			flagValues:  fixtureFlagValues(flagProjectId(nil)),
			isValid:     false,
		},
		{
			description: "sort by id",
			flagValues:  fixtureFlagValues(flagSortBy("id")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("id")),
		},
		{
			description: "sort by origin-url",
			flagValues:  fixtureFlagValues(flagSortBy("originUrl")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("originUrl")),
		},
		{
			description: "sort by status",
			flagValues:  fixtureFlagValues(flagSortBy("status")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("status")),
		},
		{
			description: "sort by created",
			flagValues:  fixtureFlagValues(flagSortBy("createdAt")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("createdAt")),
		},
		{
			description: "sort by updated",
			flagValues:  fixtureFlagValues(flagSortBy("updatedAt")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("updatedAt")),
		},
		{
			description: "sort by originUrlRelated",
			flagValues:  fixtureFlagValues(flagSortBy("originUrlRelated")),
			isValid:     true,
			expected:    fixtureInputModel(inputSortBy("originUrlRelated")),
		},
		{
			description: "invalid sort by",
			flagValues:  fixtureFlagValues(flagSortBy("invalid")),
			isValid:     false,
		},
		{
			description: "missing sort by uses default",
			flagValues: fixtureFlagValues(
				func(flagValues map[string]string) {
					delete(flagValues, sortByFlag)
				},
			),
			isValid:  true,
			expected: fixtureInputModel(inputSortBy("createdAt")),
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
			inputModel:  fixtureInputModel(inputSortBy("updatedAt")),
			expected:    fixtureRequest(requestSortBy("updatedAt")),
		},
		{
			description: "with next page id",
			inputModel:  fixtureInputModel(),
			nextPageID:  utils.Ptr(testNextPageID),
			expected:    fixtureRequest(requestNextPageID(testNextPageID)),
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

func fixtureTestResponse(mods ...func(r *testResponse)) testResponse {
	r := testResponse{
		statusCode: 200,
	}
	for _, mod := range mods {
		mod(&r)
	}
	return r
}

func responseStatus(statusCode int) func(r *testResponse) {
	return func(r *testResponse) {
		r.statusCode = statusCode
	}
}

func responseNextPageID(nextPageID *string) func(r *testResponse) {
	return func(r *testResponse) {
		r.body.NextPageIdentifier = nextPageID
	}
}

func responseDistributions(distributions ...cdn.Distribution) func(r *testResponse) {
	return func(r *testResponse) {
		r.body.Distributions = &distributions
	}
}

func fixtureDistribution(id string) cdn.Distribution {
	return cdn.Distribution{
		Id: &id,
	}
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
		limit       int
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
					responseDistributions(fixtureDistribution("dist-1")),
				),
			},
			expected: []cdn.Distribution{
				fixtureDistribution("dist-1"),
			},
		},
		{
			description: "multiple distributions, multiple pages",
			responses: []testResponse{
				fixtureTestResponse(
					responseNextPageID(utils.Ptr(testNextPageID)),
					responseDistributions(
						fixtureDistribution("dist-1"),
					),
				),
				fixtureTestResponse(
					responseDistributions(
						fixtureDistribution("dist-2"),
					),
				),
			},
			expected: []cdn.Distribution{
				fixtureDistribution("dist-1"),
				fixtureDistribution("dist-2"),
			},
		},
		{
			description: "API error",
			responses: []testResponse{
				fixtureTestResponse(
					responseStatus(500),
				),
			},
			fails: true,
		},
		{
			description: "API error on second page",
			responses: []testResponse{
				fixtureTestResponse(
					responseNextPageID(utils.Ptr(testNextPageID)),
					responseDistributions(
						fixtureDistribution("dist-1"),
					),
				),
				fixtureTestResponse(responseStatus(500)),
			},
			fails: true,
		},
		{
			description: "limit across 2 pages",
			limit:       110,
			responses: []testResponse{
				fixtureTestResponse(
					responseNextPageID(utils.Ptr(testNextPageID)),
					responseDistributions(fixtureDistributions(100)...),
				),
				fixtureTestResponse(
					responseDistributions(fixtureDistributions(10)...),
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
					m.Limit = utils.Ptr(int32(tt.limit))
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
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})

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
