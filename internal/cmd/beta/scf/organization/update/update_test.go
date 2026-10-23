package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

const (
	testRegion = "eu01"
	testName   = "my-org"
)

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &scf.APIClient{
		DefaultAPI: &scf.DefaultAPIService{},
	}
	testProjectId      = uuid.NewString()
	testOrganizationId = uuid.NewString()
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
		nameFlag:                  testName,
		suspendedFlag:             "true",
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
		Name:           utils.Ptr(testName),
		Suspended:      utils.Ptr(true),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureUpdatePayload(mods ...func(payload *scf.UpdateOrganizationPayload)) scf.UpdateOrganizationPayload {
	payload := scf.UpdateOrganizationPayload{
		Name:      utils.Ptr(testName),
		Suspended: utils.Ptr(true),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

// fixtureRequest generates an API request for tests
func fixtureRequest(mods ...func(request *scf.ApiUpdateOrganizationRequest)) scf.ApiUpdateOrganizationRequest {
	request := testClient.DefaultAPI.UpdateOrganization(testCtx, testProjectId, testRegion, testOrganizationId)
	request = request.UpdateOrganizationPayload(fixtureUpdatePayload())
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
			description: "name only",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				nameFlag:                  testName,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Suspended = nil
			}),
		},
		{
			description: "suspended only",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				suspendedFlag:             "false",
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
				model.Suspended = utils.Ptr(false)
			}),
		},
		{
			description: "no fields to update",
			argValues:   fixtureArgValues(),
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
			},
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
		expectedRequest scf.ApiUpdateOrganizationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "name only",
			model: fixtureInputModel(func(model *inputModel) {
				model.Suspended = nil
			}),
			expectedRequest: fixtureRequest(func(request *scf.ApiUpdateOrganizationRequest) {
				*request = request.UpdateOrganizationPayload(fixtureUpdatePayload(func(payload *scf.UpdateOrganizationPayload) {
					payload.Suspended = nil
				}))
			}),
		},
		{
			description: "suspended only",
			model: fixtureInputModel(func(model *inputModel) {
				model.Name = nil
			}),
			expectedRequest: fixtureRequest(func(request *scf.ApiUpdateOrganizationRequest) {
				*request = request.UpdateOrganizationPayload(fixtureUpdatePayload(func(payload *scf.UpdateOrganizationPayload) {
					payload.Name = nil
				}))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
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
	type args struct {
		model *inputModel
		resp  *scf.Organization
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
				resp:  &scf.Organization{Guid: testOrganizationId},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: &scf.Organization{Guid: testOrganizationId},
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
