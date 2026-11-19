package describe

import (
	"context"
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
	testSpaceId        = uuid.NewString()
)

func fixtureArgValues() []string {
	return []string{testSpaceId}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		organizationIdFlag:        testOrganizationId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel() *inputModel {
	return &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		OrganizationId: testOrganizationId,
		SpaceId:        testSpaceId,
	}
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
			description: "space id invalid",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
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
			description: "organization id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, organizationIdFlag)
			}),
			isValid: false,
		},
		{
			description: "organization id invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
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
	expected := testClient.DefaultAPI.GetSpace(testCtx, testProjectId, testRegion, testOrganizationId, testSpaceId)
	request := buildRequest(testCtx, fixtureInputModel(), testClient.DefaultAPI)

	diff := cmp.Diff(request, expected,
		cmp.AllowUnexported(expected),
		cmpopts.EquateComparable(testCtx),
		cmpopts.EquateComparable(testClient.DefaultAPI),
	)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestOutputResult(t *testing.T) {
	space := &scf.Space{
		Guid:      testSpaceId,
		Name:      "my-space",
		OrgId:     testOrganizationId,
		ProjectId: testProjectId,
		Region:    testRegion,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	tests := []struct {
		name         string
		space        *scf.Space
		outputFormat string
		wantErr      bool
	}{
		{
			name:  "default output",
			space: space,
		},
		{
			name:         "json output",
			space:        space,
			outputFormat: print.JSONOutputFormat,
		},
		{
			name:    "nil space",
			space:   nil,
			wantErr: true,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, tt.space); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
