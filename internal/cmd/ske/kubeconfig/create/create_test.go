package create

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	ske "github.com/stackitcloud/stackit-sdk-go/services/ske/v2api"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &ske.APIClient{DefaultAPI: &ske.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testClusterName = "cluster"

const testRegion = "eu01"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testClusterName,
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
		ClusterName: testClusterName,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *ske.ApiCreateKubeconfigRequest)) ske.ApiCreateKubeconfigRequest {
	request := testClient.DefaultAPI.CreateKubeconfig(testCtx, testProjectId, testRegion, testClusterName)
	request = request.CreateKubeconfigPayload(ske.CreateKubeconfigPayload{})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureRequestLogin() ske.ApiGetLoginKubeconfigRequest {
	return testClient.DefaultAPI.GetLoginKubeconfig(testCtx, testProjectId, testRegion, testClusterName)
}

func fixtureRequestIDP() ske.ApiGetIDPKubeconfigRequest {
	return testClient.DefaultAPI.GetIDPKubeconfig(testCtx, testProjectId, testRegion, testClusterName)
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
			description: "30d expiration time",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["expiration"] = "30d"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ExpirationTime = utils.Ptr("2592000")
			}),
		},
		{
			description: "login",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["login"] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Login = true
			}),
		},
		{
			description: "idp",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["idp"] = "true"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.IDP = true
			}),
		},
		{
			description: "custom filepath",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["filepath"] = "/path/to/config"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Filepath = utils.Ptr("/path/to/config")
			}),
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
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "disable writing and invalid output format",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[disableWritingFlag] = "true"
			}),
			isValid: false,
		},
		{
			description: "disable writing and valid output format",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[disableWritingFlag] = "true"
				flagValues[globalflags.OutputFormatFlag] = print.YAMLOutputFormat
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DisableWriting = true
				model.OutputFormat = print.YAMLOutputFormat
			}),
			isValid: true,
		},
		{
			description: "enable overwrite",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[overwriteFlag] = "true"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Overwrite = true
			}),
			isValid: true,
		},
		{
			description: "disable overwrite",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[overwriteFlag] = "false"
			}),
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Overwrite = false
			}),
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequestCreate(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest ske.ApiCreateKubeconfigRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "expiration time",
			model: fixtureInputModel(func(model *inputModel) {
				model.ExpirationTime = utils.Ptr("2592000")
			}),
			expectedRequest: fixtureRequest().CreateKubeconfigPayload(ske.CreateKubeconfigPayload{
				ExpirationSeconds: utils.Ptr("2592000")}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, _ := buildRequestCreate(testCtx, tt.model, testClient)
			assertNoDiff(t, request, tt.expectedRequest)
		})
	}
}

func assertNoDiff(t *testing.T, actual, expected any) {
	t.Helper()
	diff := cmp.Diff(actual, expected,
		cmp.AllowUnexported(expected),
		cmpopts.EquateComparable(testCtx),
		cmpopts.EquateComparable(testClient.DefaultAPI),
	)
	if diff != "" {
		t.Fatalf("Data does not match: %s", diff)
	}
}

func TestBuildRequestLogin(t *testing.T) {
	model := fixtureInputModel()
	expectedRequest := fixtureRequestLogin()
	request, _ := buildRequestLogin(testCtx, model, testClient)
	assertNoDiff(t, request, expectedRequest)
}

func TestBuildRequestIDP(t *testing.T) {
	model := fixtureInputModel()
	expectedRequest := fixtureRequestIDP()
	request, _ := buildRequestIDP(testCtx, model, testClient)
	assertNoDiff(t, request, expectedRequest)
}

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat   string
		clusterName    string
		kubeconfigPath string
		respKubeconfig *ske.Kubeconfig
		respLogin      *ske.LoginKubeconfig
		respIDP        *ske.IDPKubeconfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
		{
			name: "missing kubeconfig",
			args: args{
				respLogin: &ske.LoginKubeconfig{},
			},
			wantErr: false,
		},
		{
			name: "missing login",
			args: args{
				respKubeconfig: &ske.Kubeconfig{},
			},
			wantErr: false,
		},
		{
			name: "missing idp",
			args: args{
				respIDP: &ske.IDPKubeconfig{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.clusterName, tt.args.kubeconfigPath, tt.args.respKubeconfig, tt.args.respLogin, tt.args.respIDP); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
