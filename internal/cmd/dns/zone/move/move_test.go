package move

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	dns "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

const (
	testDnsName                 = "www.example.com"
	testCode                    = "some-move-code"
	testCodeArg                 = "@" + testCodeFile
	testCodeFile                = "move-code.txt"
	testCodeFileTrailingNewline = "move-code-newline.txt"
)

func cmdFactory(params *types.CmdParams) *cobra.Command {
	params.Fs = fstest.MapFS{
		testCodeFile:                {Data: []byte(testCode)},
		testCodeFileTrailingNewline: {Data: []byte(testCode + "\n")},
	}
	return NewCmd(params)
}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &dns.APIClient{DefaultAPI: &dns.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testTargetProjectId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		zoneDomainNameFlag:        testDnsName,
		targetProjectIdFlag:       testTargetProjectId,
		codeFlag:                  testCodeArg,
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
		},
		ZoneDomainName:  testDnsName,
		TargetProjectId: testTargetProjectId,
		Code:            utils.Ptr(testCode),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *dns.ApiMoveZoneRequest)) dns.ApiMoveZoneRequest {
	request := testClient.DefaultAPI.MoveZone(testCtx, testTargetProjectId)
	request = request.MoveZonePayload(dns.MoveZonePayload{
		Code:        testCode,
		ZoneDnsName: testDnsName,
	})
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
			description: "without code",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, codeFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Code = nil
			}),
		},
		{
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "zone domain name missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, zoneDomainNameFlag)
			}),
			isValid: false,
		},
		{
			description: "code file with trailing newline",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[codeFlag] = "@" + testCodeFileTrailingNewline
			}),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "target project id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, targetProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "target project id missing without code",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, targetProjectIdFlag)
				delete(flagValues, codeFlag)
			}),
			isValid: false,
		},
		{
			description: "target project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[targetProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "project id missing but code given",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id and code missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
				delete(flagValues, codeFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "literal code rejected",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[codeFlag] = testCode
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, cmdFactory, parseInput, tt.expectedModel, nil, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		code            string
		expectedRequest dns.ApiMoveZoneRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			code:            testCode,
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, tt.code, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, dns.DefaultAPIService{}),
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
		resp  *dns.Message
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
			name: "base",
			args: args{
				model: fixtureInputModel(),
				resp:  &dns.Message{},
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
