package list

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

var projectIdFlag = globalflags.ProjectIdFlag.FlagName()

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()
var testZoneId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:   testProjectId,
		zoneIdFlag:      testZoneId,
		nameLikeFlag:    "some-pattern",
		activeFlag:      "true",
		orderByNameFlag: "asc",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureFlagModel(mods ...func(model *flagModel)) *flagModel {
	model := &flagModel{
		ProjectId:   testProjectId,
		ZoneId:      testZoneId,
		NameLike:    utils.Ptr("some-pattern"),
		Active:      utils.Ptr(true),
		OrderByName: utils.Ptr("asc"),
		PageSize:    100,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiGetRecordSetsRequest)) dns.ApiGetRecordSetsRequest {
	request := testClient.GetRecordSets(testCtx, testProjectId, testZoneId)
	request = request.NameLike("some-pattern")
	request = request.ActiveEq(true)
	request = request.OrderByName("ASC")
	request = request.PageSize(100)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		description   string
		flagValues    map[string]string
		isValid       bool
		expectedModel *flagModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureFlagModel(),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				projectIdFlag: testProjectId,
				zoneIdFlag:    testZoneId,
			},
			isValid: true,
			expectedModel: &flagModel{
				ProjectId: testProjectId,
				ZoneId:    testZoneId,
				PageSize:  100, // Default value
			},
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
			description: "name like empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameLikeFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.NameLike = utils.Ptr("")
			}),
		},
		{
			description: "is active = false",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "false"
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.Active = utils.Ptr(false)
			}),
		},
		{
			description: "is active invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "is active invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "invalid"
			}),
			isValid: false,
		},
		{
			description: "order by name desc",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = "desc"
			}),
			isValid: true,
			expectedModel: fixtureFlagModel(func(model *flagModel) {
				model.OrderByName = utils.Ptr("desc")
			}),
		},
		{
			description: "order by name invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "order by name invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = "invalid"
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
			cmd := &cobra.Command{}
			configureFlags(cmd)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseFlags(cmd)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *flagModel
		page            int
		expectedRequest dns.ApiGetRecordSetsRequest
	}{
		{
			description:     "base",
			model:           fixtureFlagModel(),
			page:            1,
			expectedRequest: fixtureRequest().Page(1),
		},
		{
			description:     "base 2",
			model:           fixtureFlagModel(),
			page:            10,
			expectedRequest: fixtureRequest().Page(10),
		},
		{
			description: "required fields only",
			model: &flagModel{
				ProjectId: testProjectId,
				ZoneId:    testZoneId,
				PageSize:  10,
			},
			page:            1,
			expectedRequest: testClient.GetRecordSets(testCtx, testProjectId, testZoneId).Page(1).PageSize(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, tt.page)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestFetchRecordSets(t *testing.T) {
	tests := []struct {
		description         string
		model               *flagModel
		totalItems          int
		apiCallFails        bool
		expectedNumAPICalls int
		expectedNumItems    int
	}{
		{
			description:         "no limit and pageSize>totalItems",
			model:               fixtureFlagModel(),
			totalItems:          10,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description:         "no limit and pageSize<totalItems",
			model:               fixtureFlagModel(),
			totalItems:          320,
			expectedNumAPICalls: 4,
			apiCallFails:        false,
			expectedNumItems:    320,
		},
		{
			description:         "no limit and pageSize<totalItems 2",
			model:               fixtureFlagModel(),
			totalItems:          200,
			expectedNumAPICalls: 3, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    200,
		},
		{
			description:         "no limit and pageSize=totalItems",
			model:               fixtureFlagModel(),
			totalItems:          100,
			expectedNumAPICalls: 2, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    100,
		},
		{
			description: "limit<pageSize",
			model: fixtureFlagModel(func(model *flagModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
			totalItems:          100,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description: "limit>totalItems and pageSize>totalItems",
			model: fixtureFlagModel(func(model *flagModel) {
				model.Limit = utils.Ptr(int64(200))
				model.PageSize = 300
			}),
			totalItems:          50,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    50,
		},
		{
			description: "limit>totalItems and pageSize<totalItems",
			model: fixtureFlagModel(func(model *flagModel) {
				model.Limit = utils.Ptr(int64(200))
				model.PageSize = 30
			}),
			totalItems:          50,
			expectedNumAPICalls: 2,
			apiCallFails:        false,
			expectedNumItems:    50,
		},
		{
			description:  "request fails",
			model:        fixtureFlagModel(),
			totalItems:   100,
			apiCallFails: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			numAPICalls := 0
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				numAPICalls++

				w.Header().Set("Content-Type", "application/json")
				if tt.apiCallFails {
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write([]byte("{\"message\": \"Something bad happened\""))
					if err != nil {
						t.Errorf("Failed to write bad response: %v", err)
					}
					return
				}

				query := r.URL.Query()
				pageStr := query.Get("page")
				if pageStr == "" {
					t.Errorf("Expected query param page to be set")
				}
				page, err := strconv.Atoi(pageStr)
				if err != nil {
					t.Errorf("Failed to parse query param page: %v", err)
				}
				pageSizeStr := query.Get("pageSize")
				if pageSizeStr == "" {
					t.Errorf("Expected query param pageSize to be set")
				}
				pageSize, err := strconv.Atoi(pageSizeStr)
				if err != nil {
					t.Errorf("Failed to parse query param pageSize: %v", err)
				}

				offset := (page - 1) * pageSize

				var numItemsToReturn int
				if offset >= tt.totalItems {
					numItemsToReturn = 0 // Total items reached
				} else if offset+pageSize < tt.totalItems {
					numItemsToReturn = pageSize // Full intermediate page
				} else {
					numItemsToReturn = tt.totalItems - offset // Last page
				}

				recordSets := make([]dns.RecordSet, numItemsToReturn)
				mockedResp := dns.RecordSetsResponse{
					RrSets: &recordSets,
				}

				mockedRespBytes, err := json.Marshal(mockedResp)
				if err != nil {
					t.Fatalf("Failed to marshal mocked response: %v", err)
				}

				_, err = w.Write(mockedRespBytes)
				if err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			})
			mockedServer := httptest.NewServer(handler)
			defer mockedServer.Close()
			client, err := dns.NewAPIClient(
				sdkConfig.WithEndpoint(mockedServer.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("Failed to initialize client: %v", err)
			}

			recordSets, err := fetchRecordSets(testCtx, tt.model, client)
			if err != nil {
				if !tt.apiCallFails {
					t.Fatalf("did not fail on invalid input")
				}
				return
			}
			if err == nil && tt.apiCallFails {
				t.Fatalf("did not fail on invalid input")
			}
			if numAPICalls != tt.expectedNumAPICalls {
				t.Fatalf("Expected %d API calls, got %d", tt.expectedNumAPICalls, numAPICalls)
			}
			if len(recordSets) != tt.expectedNumItems {
				t.Fatalf("Expected %d recordSets, got %d", tt.totalItems, len(recordSets))
			}
		})
	}
}
