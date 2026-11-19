package update

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	git "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

type testCtxKey struct{}

var (
	testCtx        = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient     = &git.APIClient{DefaultAPI: &git.DefaultAPIService{}}
	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		aclFlag:                   "0.0.0.0/0",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{ProjectId: testProjectId, Verbosity: globalflags.VerbosityDefault},
		InstanceId:      testInstanceId,
		Acl:             []string{"0.0.0.0/0"},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *git.ApiPatchInstanceRequest)) git.ApiPatchInstanceRequest {
	request := testClient.DefaultAPI.PatchInstance(testCtx, testProjectId, testInstanceId)
	request = request.PatchInstancePayload(git.PatchInstancePayload{
		Acl: []string{"0.0.0.0/0"},
	})
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
			description: "multiple acl entries",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = "1.2.3.0/24,4.5.6.0/24"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Acl = []string{"1.2.3.0/24", "4.5.6.0/24"}
			}),
		},
		{
			description: "invalid acl",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[aclFlag] = "invalid-cidr"
			}),
			isValid: false,
		},
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no values to update",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			isValid: false,
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
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
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
		expectedRequest git.ApiPatchInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "multiple acl entries",
			model: fixtureInputModel(func(model *inputModel) {
				model.Acl = []string{"1.2.3.0/24", "4.5.6.0/24"}
			}),
			expectedRequest: fixtureRequest(func(request *git.ApiPatchInstanceRequest) {
				*request = request.PatchInstancePayload(git.PatchInstancePayload{
					Acl: []string{"1.2.3.0/24", "4.5.6.0/24"},
				})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)
			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, git.DefaultAPIService{}),
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
		model *inputModel
		resp  *git.Instance
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{},
			wantErr: true,
		},
		{
			name: "nil model",
			args: args{
				resp: &git.Instance{},
			},
			wantErr: true,
		},
		{
			name: "nil resp",
			args: args{
				model: fixtureInputModel(),
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				model: fixtureInputModel(),
				resp:  &git.Instance{},
			},
			wantErr: false,
		},
		{
			name: "async",
			args: args{
				model: fixtureInputModel(func(model *inputModel) { model.Async = true }),
				resp:  &git.Instance{},
			},
			wantErr: false,
		},
		{
			name: "json output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) { model.OutputFormat = print.JSONOutputFormat }),
				resp:  &git.Instance{},
			},
			wantErr: false,
		},
		{
			name: "yaml output",
			args: args{
				model: fixtureInputModel(func(model *inputModel) { model.OutputFormat = print.YAMLOutputFormat }),
				resp:  &git.Instance{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, "label", tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
