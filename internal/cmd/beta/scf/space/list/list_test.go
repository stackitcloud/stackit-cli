package list

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const testRegion = "eu01"

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &scf.APIClient{
		DefaultAPI: &scf.DefaultAPIService{},
	}
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		organizationIdFlag:        testOrganizationId,
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
		OrganizationId: testOrganizationId,
		Limit:          new(int64(10)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "limit invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "organization id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
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
		expectedRequest scf.ApiListSpacesRequest
	}{
		{
			description:     "limit below page size",
			model:           fixtureInputModel(),
			expectedRequest: testClient.DefaultAPI.ListSpaces(testCtx, testProjectId, testRegion, testOrganizationId).Page(1).PerPage(10),
		},
		{
			description: "no limit",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = nil
			}),
			expectedRequest: testClient.DefaultAPI.ListSpaces(testCtx, testProjectId, testRegion, testOrganizationId).Page(1).PerPage(maxPageSize),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI, 1, maxPageSize)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.EquateComparable(testClient.DefaultAPI),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	resp := &scf.SpacesList{
		Pagination: scf.Pagination{TotalPages: new(int64(2))},
		Resources:  []scf.Space{{Guid: uuid.NewString(), Name: "my-space", CreatedAt: time.Now()}},
	}
	tests := []struct {
		name         string
		outputFormat string
		limit        *int64
		resp         *scf.SpacesList
		expectWarn   bool
	}{
		{
			name:       "default output",
			resp:       resp,
			expectWarn: true,
		},
		{
			name:         "json output",
			outputFormat: print.JSONOutputFormat,
			resp:         resp,
		},
		{
			name:  "limit reached",
			limit: new(int64(1)),
			resp:  resp,
		},
		{
			name:       "limit above page size",
			limit:      new(int64(maxPageSize + 1)),
			resp:       resp,
			expectWarn: true,
		},
		{
			name: "empty",
			resp: &scf.SpacesList{},
		},
		{
			name: "nil response",
			resp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(params.Printer, tt.outputFormat, tt.limit, "org", tt.resp); err != nil {
				t.Errorf("outputResult() error = %v", err)
			}
			warned := strings.Contains(params.Err.String(), "more than 1 page")
			if warned != tt.expectWarn {
				t.Errorf("expected warning = %t, got %t (stderr: %q)", tt.expectWarn, warned, params.Err.String())
			}
		})
	}
}
