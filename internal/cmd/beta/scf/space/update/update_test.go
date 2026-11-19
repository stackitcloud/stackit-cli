package update

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
	testName   = "new-name"
)

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
		nameFlag:                  testName,
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
		SpaceId:        testSpaceId,
		Name:           testName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "name missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
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
	expected := testClient.DefaultAPI.UpdateSpace(testCtx, testProjectId, testRegion, testOrganizationId, testSpaceId).
		UpdateSpacePayload(scf.UpdateSpacePayload{Name: testName})
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
	type args struct {
		model *inputModel
		resp  *scf.Space
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default output",
			args: args{
				model: fixtureInputModel(),
				resp:  &scf.Space{Guid: testSpaceId},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: &scf.Space{Guid: testSpaceId},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
