package list

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &scf.APIClient{DefaultAPI: &scf.DefaultAPIService{}}
	testProjectId    = uuid.NewString()
	testOrganization = uuid.NewString()
)

const (
	testLimit  int64 = 10
	testRegion       = "eu01"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		limitFlag:                 strconv.FormatInt(testLimit, 10),
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
		Limit: new(testLimit),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *scf.ApiListOrganizationsRequest)) scf.ApiListOrganizationsRequest {
	request := testClient.DefaultAPI.ListOrganizations(testCtx, testProjectId, testRegion).WithDetails(true).Page(1).PerPage(testLimit)
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
			description:   "no values",
			flagValues:    map[string]string{},
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "missing project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid:       false,
			expectedModel: nil,
		},
		{
			description: "no limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, limitFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = nil
			}),
		},
		{
			description: "invalid limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid:       false,
			expectedModel: nil,
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
		isValid         bool
		expectedRequest scf.ApiListOrganizationsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI, 1, testLimit)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, scf.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func fixtureResponse(count int) []scf.OrganizationsListItem {
	return make([]scf.OrganizationsListItem, count)
}

func TestFetchOrganizations(t *testing.T) {
	tests := []struct {
		description   string
		limit         *int64
		maxPageSize   int
		response      []scf.OrganizationsListItem
		wantErr       bool
		expectedCount int
	}{
		{
			description:   "no organizations",
			maxPageSize:   10,
			response:      fixtureResponse(0),
			expectedCount: 0,
		},
		{
			description:   "Limit is smaller than default page size",
			limit:         new(int64(2)),
			maxPageSize:   10,
			expectedCount: 2,
			response:      fixtureResponse(3),
		},
		{
			description:   "Limit is unspecified",
			maxPageSize:   2,
			expectedCount: 2,
			response:      fixtureResponse(3),
		},
		{
			description: "page size smaller than 1 should error",
			maxPageSize: 0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				resources := tt.response
				if tt.maxPageSize < len(tt.response) {
					resources = resources[:tt.maxPageSize]
				}
				serverResponse := scf.OrganizationsList{
					Resources: resources,
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(serverResponse); err != nil {
					t.Fatalf("write: %v", err)
				}
			})
			server := httptest.NewServer(handler)
			defer server.Close()
			client, err := scf.NewAPIClient(
				sdkConfig.WithEndpoint(server.URL),
				sdkConfig.WithoutAuthentication(),
			)
			if err != nil {
				t.Fatalf("failed to create test client: %v", err)
			}

			model := fixtureInputModel(func(model *inputModel) {
				model.Limit = tt.limit
			})
			got, err := fetchOrganizations(testCtx, nil, model, tt.maxPageSize, client.DefaultAPI)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("fetchOrganizations() expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("fetchOrganizations() unexpected error: %v", err)
			}
			if got == nil {
				t.Errorf("fetchOrganizations() returned nil slice, expected empty slice")
			}
			if len(got) != tt.expectedCount {
				t.Errorf("fetchOrganizations() expected %d organizations, got %d", tt.expectedCount, len(got))
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		projectLabel  string
		organizations []scf.OrganizationsListItem
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
			name: "nil organizations slice",
			args: args{
				projectLabel:  testProjectId,
				organizations: nil,
			},
			wantErr: false,
		},
		{
			name: "empty organizations slice",
			args: args{
				projectLabel:  testProjectId,
				organizations: []scf.OrganizationsListItem{},
			},
			wantErr: false,
		},
		{
			name: "empty organization struct",
			args: args{
				projectLabel: testProjectId,
				organizations: []scf.OrganizationsListItem{
					{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid organizations",
			args: args{
				projectLabel: testProjectId,
				organizations: []scf.OrganizationsListItem{
					{
						Guid:      testOrganization,
						Name:      new("my-org"),
						Status:    "active",
						CreatedAt: new(time.Now()),
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.organizations); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
