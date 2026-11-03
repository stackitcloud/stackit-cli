package add

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	loadbalancer "github.com/stackitcloud/stackit-sdk-go/services/loadbalancer/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &loadbalancer.APIClient{DefaultAPI: &loadbalancer.DefaultAPIService{}}
var testProjectId = uuid.NewString()

const (
	testPwFile = "pw.txt"
	testPw     = "file-pw"
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		displayNameFlag:           "name",
		usernameFlag:              "username",
		passwordFlag:              "@" + testPwFile,
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
		DisplayName: utils.Ptr("name"),
		Username:    utils.Ptr("username"),
		Password:    utils.Ptr(testPw),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *loadbalancer.ApiCreateCredentialsRequest)) loadbalancer.ApiCreateCredentialsRequest {
	request := testClient.DefaultAPI.CreateCredentials(testCtx, testProjectId, testRegion)
	request = request.CreateCredentialsPayload(loadbalancer.CreateCredentialsPayload{
		DisplayName: utils.Ptr("name"),
		Username:    utils.Ptr("username"),
		Password:    new(testPw),
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		stdin         string
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
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "display name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "username name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, usernameFlag)
			}),
			isValid: false,
		},
		{
			description: "password from explicit stdin",
			stdin:       "stdin-pw\n",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[passwordFlag] = "-"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Password = new("stdin-pw")
			}),
		},
		{
			// Deprecated compatibility behavior; remove this case when automatic prompting is removed in 2027-04 (STACKITCLI-549).
			description: "password from implicit stdin",
			stdin:       "stdin-pw\n",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, passwordFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Password = new("stdin-pw")
			}),
		},
		{
			description: "literal password rejected",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[passwordFlag] = "top-secret"
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, func(params *types.CmdParams) *cobra.Command {
				params.Fs = fstest.MapFS{
					testPwFile: {Data: []byte(testPw)},
				}
				params.Printer.StdIn = strings.NewReader(tt.stdin)
				return NewCmd(params)
			}, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest loadbalancer.ApiCreateCredentialsRequest
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
				cmpopts.EquateComparable(testCtx, loadbalancer.DefaultAPIService{}),
				cmpopts.IgnoreFields(loadbalancer.ApiCreateCredentialsRequest{}, "xRequestID"),
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
		resp         *loadbalancer.CreateCredentialsResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "credentials response with empty Credential",
			args: args{
				resp: &loadbalancer.CreateCredentialsResponse{},
			},
			wantErr: true,
		},
		{
			name: "only credentials response as argument",
			args: args{
				resp: &loadbalancer.CreateCredentialsResponse{Credential: &loadbalancer.CredentialsResponse{}},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
