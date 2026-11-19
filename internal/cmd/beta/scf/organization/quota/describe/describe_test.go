package describe

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
)

var (
	testCtx            = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient         = &scf.DefaultAPIService{}
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
	testQuotaId        = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testOrganizationId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		quotaIdFlag:               testQuotaId,
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
		QuotaId:        new(testQuotaId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *scf.ApiGetOrganizationQuotaRequest)) scf.ApiGetOrganizationQuotaRequest {
	request := testClient.GetOrganizationQuota(testCtx, testProjectId, testRegion, testOrganizationId, testQuotaId)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "organization id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "quota id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, quotaIdFlag)
			}),
			isValid: false,
		},
		{
			description: "quota id invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[quotaIdFlag] = "invalid-uuid"
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
		expectedRequest scf.ApiGetOrganizationQuotaRequest
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
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.EquateComparable(testClient),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name         string
		quota        *scf.Quota
		outputFormat string
		wantErr      bool
	}{
		{
			name:    "default output",
			quota:   fixtureQuota(),
			wantErr: false,
		},
		{
			name:         "json output",
			quota:        fixtureQuota(),
			outputFormat: print.JSONOutputFormat,
			wantErr:      false,
		},
		{
			name:    "nil quota",
			quota:   nil,
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.quota); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func fixtureQuota() *scf.Quota {
	return &scf.Quota{
		Guid:       testQuotaId,
		Name:       "my-quota",
		OrgId:      new(testOrganizationId),
		PlatformId: uuid.NewString(),
		ProjectId:  testProjectId,
		Region:     testRegion,
		Apps: scf.QuotaApps{
			TotalInstances:               *scf.NewNullableInt64(new(int64(10))),
			TotalMemoryInMb:              *scf.NewNullableInt64(new(int64(1024))),
			PerProcessMemoryInMb:         *scf.NewNullableInt64(nil),
			PerAppTasks:                  *scf.NewNullableInt64(new(int64(5))),
			LogRateLimitInBytesPerSecond: *scf.NewNullableInt64(nil),
		},
		Domains: scf.QuotaDomains{
			TotalDomains: *scf.NewNullableInt64(new(int64(2))),
		},
		Routes: scf.QuotaRoutes{
			TotalRoutes:        *scf.NewNullableInt64(new(int64(20))),
			TotalReservedPorts: *scf.NewNullableInt64(nil),
		},
		Services: scf.QuotaServices{
			PaidServicesAllowed:   true,
			TotalServiceInstances: *scf.NewNullableInt64(new(int64(10))),
			TotalServiceKeys:      *scf.NewNullableInt64(new(int64(10))),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
