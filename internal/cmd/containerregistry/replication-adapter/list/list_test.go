package list

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	projectIdFlag = globalflags.ProjectIdFlag
	testRegion    = "eu01"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &registry.APIClient{DefaultAPI: &registry.DefaultAPIService{}}
	testProjectId = uuid.NewString()
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:          testProjectId,
		globalflags.RegionFlag: testRegion,
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
			Region:    testRegion,
		},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *registry.ApiListReplicationAdaptersRequest)) registry.ApiListReplicationAdaptersRequest {
	request := testClient.DefaultAPI.ListReplicationAdapters(testCtx, testProjectId, testRegion)
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
			description: "with limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "10"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(10))
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "limit invalid: zero",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
			}),
			isValid: false,
		},
		{
			description: "limit invalid: negative",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "-1"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, []string{}, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest registry.ApiListReplicationAdaptersRequest
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
				cmp.AllowUnexported(tt.expectedRequest, registry.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat string
		projectLabel string
		adapters     []registry.ReplicationAdapter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty adapters",
			args: args{
				projectLabel: "my-project",
				adapters:     []registry.ReplicationAdapter{},
			},
			wantErr: false,
		},
		{
			name: "single item",
			args: args{
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple items",
			args: args{
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
					{
						RegistryType:   "harbor",
						DisplayName:    "Harbor",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "json output format",
			args: args{
				outputFormat: print.JSONOutputFormat,
				projectLabel: "my-project",
				adapters: []registry.ReplicationAdapter{
					{
						RegistryType:   "docker-hub",
						DisplayName:    "Docker Hub",
						CredentialType: registry.REPLICATIONADAPTERCREDENTIALTYPE_BASIC,
					},
				},
			},
			wantErr: false,
		},
	}

	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.adapters)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
