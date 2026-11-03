package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
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
	testProjectId  = uuid.NewString()
	testPlatformId = uuid.NewString()
)

// fixtureFlagValues generates a map of flag values for tests
func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		nameFlag:                  testName,
		platformIdFlag:            testPlatformId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

// fixtureInputModel generates an input model for tests
func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		Name:       new(testName),
		PlatformId: new(testPlatformId),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

// fixtureCreatePayload generates a CreateOrganizationPayload for tests
func fixtureCreatePayload(mods ...func(payload *scf.CreateOrganizationPayload)) scf.CreateOrganizationPayload {
	payload := scf.CreateOrganizationPayload{
		Name:       testName,
		PlatformId: new(testPlatformId),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

// fixtureRequest generates an API request for tests
func fixtureRequest(mods ...func(request *scf.ApiCreateOrganizationRequest)) scf.ApiCreateOrganizationRequest {
	request := testClient.DefaultAPI.CreateOrganization(testCtx, testProjectId, testRegion)
	request = request.CreateOrganizationPayload(fixtureCreatePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
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
		{
			description: "platform id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[platformIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "required fields only",
			flagValues: map[string]string{
				globalflags.ProjectIdFlag: testProjectId,
				globalflags.RegionFlag:    testRegion,
				nameFlag:                  testName,
			},
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.PlatformId = nil
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, func(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
				return parseInput(p, cmd)
			}, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest scf.ApiCreateOrganizationRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "no platform id",
			model: fixtureInputModel(func(model *inputModel) {
				model.PlatformId = nil
			}),
			expectedRequest: fixtureRequest(func(request *scf.ApiCreateOrganizationRequest) {
				*request = request.CreateOrganizationPayload(fixtureCreatePayload(func(payload *scf.CreateOrganizationPayload) {
					payload.PlatformId = nil
				}))
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient.DefaultAPI)
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
		model        *inputModel
		projectLabel string
		resp         *scf.OrganizationCreateResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default output",
			args: args{
				model:        fixtureInputModel(),
				projectLabel: "my-project",
				resp:         &scf.OrganizationCreateResponse{Guid: "org-id-123"},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: &scf.OrganizationCreateResponse{Guid: "org-id-123"},
			},
			wantErr: false,
		},
		{
			name: "nil response - default output",
			args: args{
				model: fixtureInputModel(),
				resp:  nil,
			},
			wantErr: false,
		},
		{
			name: "nil response - json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) {
					model.OutputFormat = print.JSONOutputFormat
				}),
				resp: nil,
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
