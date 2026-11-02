package create

import (
	"context"
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

const (
	testRegion = "eu01"
	testName   = "my-space"
)

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
		nameFlag:                  testName,
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
		Name:           testName,
	}
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
		{
			description: "organization id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[organizationIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
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
	expected := testClient.DefaultAPI.CreateSpace(testCtx, testProjectId, testRegion, testOrganizationId).
		CreateSpacePayload(scf.CreateSpacePayload{Name: testName})
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
	tests := []struct {
		name         string
		outputFormat string
		resp         *scf.Space
	}{
		{
			name: "default output",
			resp: &scf.Space{Guid: uuid.NewString()},
		},
		{
			name:         "json output",
			outputFormat: print.JSONOutputFormat,
			resp:         &scf.Space{Guid: uuid.NewString()},
		},
		{
			name: "nil response",
			resp: nil,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, "org", tt.resp); err != nil {
				t.Errorf("outputResult() error = %v", err)
			}
		})
	}
}
