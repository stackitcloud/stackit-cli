package list

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/dns"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{}
var testProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		nameLikeFlag:              "some-pattern",
		orderByNameFlag:           "asc",
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
		},
		NameLike:    utils.Ptr("some-pattern"),
		OrderByName: utils.Ptr("asc"),
		PageSize:    pageSizeDefault,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiListZonesRequest)) dns.ApiListZonesRequest {
	request := testClient.ListZones(testCtx, testProjectId)
	request = request.NameLike("some-pattern")
	request = request.OrderByName("ASC")
	request = request.PageSize(pageSizeDefault)
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
			description: "include deleted zones",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[includeDeletedFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IncludeDeleted = true
			}),
		},
		{
			description: "active zones",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Active = true
			}),
		},
		{
			description: "inactive zones",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[inactiveFlag] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Inactive = true
			}),
		},
		{
			description: "active and inactive zones",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[activeFlag] = "true"
				flagValues[inactiveFlag] = "true"
			}),
			isValid: false,
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
			},
			isValid: true,
			expectedModel: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				PageSize: pageSizeDefault,
			},
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
			description: "name like empty",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[nameLikeFlag] = ""
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NameLike = utils.Ptr("")
			}),
		},
		{
			description: "order by name desc",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[orderByNameFlag] = "desc"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
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
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		page            int
		expectedRequest dns.ApiListZonesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			page:            1,
			expectedRequest: fixtureRequest().StateNeq(deleteSucceededState).Page(1),
		},
		{
			description:     "base 2",
			model:           fixtureInputModel(),
			page:            10,
			expectedRequest: fixtureRequest().StateNeq(deleteSucceededState).Page(10),
		},
		{
			description: "include deleted zones",
			model: fixtureInputModel(func(model *inputModel) {
				model.IncludeDeleted = true
			}),
			page:            1,
			expectedRequest: fixtureRequest().Page(1),
		},
		{
			description: "active zones",
			model: fixtureInputModel(func(model *inputModel) {
				model.Active = true
			}),
			page:            1,
			expectedRequest: fixtureRequest().ActiveEq(true).StateNeq(deleteSucceededState).Page(1),
		},
		{
			description: "inactive zones",
			model: fixtureInputModel(func(model *inputModel) {
				model.Inactive = true
			}),
			page:            1,
			expectedRequest: fixtureRequest().ActiveEq(false).StateNeq(deleteSucceededState).Page(1),
		},
		{
			description: "required fields only",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
					Verbosity: globalflags.VerbosityDefault,
				},
				PageSize: pageSizeDefault,
			},
			page:            1,
			expectedRequest: testClient.ListZones(testCtx, testProjectId).Page(1).PageSize(pageSizeDefault).StateNeq(deleteSucceededState),
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

func TestFetchZones(t *testing.T) {
	tests := []struct {
		description         string
		model               *inputModel
		totalItems          int
		apiCallFails        bool
		expectedNumAPICalls int
		expectedNumItems    int
	}{
		{
			description:         "no limit and pageSize>totalItems",
			model:               fixtureInputModel(),
			totalItems:          10,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description:         "no limit and pageSize<totalItems",
			model:               fixtureInputModel(),
			totalItems:          320,
			expectedNumAPICalls: 4,
			apiCallFails:        false,
			expectedNumItems:    320,
		},
		{
			description:         "no limit and pageSize<totalItems 2",
			model:               fixtureInputModel(),
			totalItems:          200,
			expectedNumAPICalls: 3, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    200,
		},
		{
			description:         "no limit and pageSize=totalItems",
			model:               fixtureInputModel(),
			totalItems:          100,
			expectedNumAPICalls: 2, // Last call will return no items
			apiCallFails:        false,
			expectedNumItems:    100,
		},
		{
			description: "limit<pageSize",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
			totalItems:          100,
			expectedNumAPICalls: 1,
			apiCallFails:        false,
			expectedNumItems:    10,
		},
		{
			description: "limit>totalItems and pageSize>totalItems",
			model: fixtureInputModel(func(model *inputModel) {
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
			model: fixtureInputModel(func(model *inputModel) {
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
			model:        fixtureInputModel(),
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

				zones := make([]dns.Zone, numItemsToReturn)
				mockedResp := dns.ListZonesResponse{
					Zones: &zones,
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

			zones, err := fetchZones(testCtx, tt.model, client)
			if err != nil {
				if !tt.apiCallFails {
					t.Fatalf("did not fail on invalid input")
				}
				return
			}
			if tt.apiCallFails {
				t.Fatalf("did not fail on invalid input")
			}
			if numAPICalls != tt.expectedNumAPICalls {
				t.Fatalf("Expected %d API calls, got %d", tt.expectedNumAPICalls, numAPICalls)
			}
			if len(zones) != tt.expectedNumItems {
				t.Fatalf("Expected %d zones, got %d", tt.totalItems, len(zones))
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		zones        []dns.Zone
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
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.zones); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
