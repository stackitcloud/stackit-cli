package list

import (
	"context"
	"strconv"
	"strings"
	"testing"

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

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &scf.APIClient{DefaultAPI: &scf.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testPlatformId = uuid.NewString()
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

func fixtureRequest(mods ...func(request *scf.ApiListPlatformsRequest)) scf.ApiListPlatformsRequest {
	request := testClient.DefaultAPI.ListPlatforms(testCtx, testProjectId, testRegion).Page(1).PerPage(testLimit)
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
		expectedRequest scf.ApiListPlatformsRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no limit set",
			model: fixtureInputModel(func(model *inputModel) {
				model.Limit = nil
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *scf.ApiListPlatformsRequest) {
				*request = request.PerPage(maxPageSize)
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI)

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

func TestOutputResult(t *testing.T) {
	resp := &scf.PlatformList{
		Pagination: scf.Pagination{TotalPages: new(int64(2))},
		Resources: []scf.Platforms{
			{
				Guid:        testPlatformId,
				DisplayName: "my-platform",
				Region:      testRegion,
				SystemId:    "my-system",
				ApiUrl:      "https://api.example.com",
			},
		},
	}
	tests := []struct {
		name         string
		outputFormat string
		limit        *int64
		resp         *scf.PlatformList
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
			resp: &scf.PlatformList{},
		},
		{
			name: "nil response",
			resp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := testparams.NewTestParams()
			if err := outputResult(params.Printer, tt.outputFormat, testProjectId, tt.limit, tt.resp); err != nil {
				t.Errorf("outputResult() error = %v", err)
			}
			warned := strings.Contains(params.Err.String(), "more than 1 page")
			if warned != tt.expectWarn {
				t.Errorf("expected warning = %t, got %t (stderr: %q)", tt.expectWarn, warned, params.Err.String())
			}
		})
	}
}
