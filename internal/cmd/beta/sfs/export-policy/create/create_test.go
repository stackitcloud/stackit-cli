package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

var projectIdFlag = globalflags.ProjectIdFlag
var regionFlag = globalflags.RegionFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sfs.APIClient{}

var testProjectId = uuid.NewString()
var testRegion = "eu01"
var testName = "test-name"
var testRulesString = "[]"
var testRules = &[]sfs.CreateShareExportPolicyRequestRule{}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,

		nameFlag:  testName,
		rulesFlag: testRulesString,
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
		Name:  testName,
		Rules: testRules,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiCreateShareExportPolicyRequest)) sfs.ApiCreateShareExportPolicyRequest {
	request := testClient.CreateShareExportPolicy(testCtx, testProjectId, testRegion)
	request = request.CreateShareExportPolicyPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *sfs.CreateShareExportPolicyPayload)) sfs.CreateShareExportPolicyPayload {
	payload := sfs.CreateShareExportPolicyPayload{
		Name:  utils.Ptr(testName),
		Rules: testRules,
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
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
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "required only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, rulesFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Rules = nil
			}),
		},
		{
			description: "required read rules from file",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[rulesFlag] = "@../test-files/rules-example.json"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Rules = &[]sfs.CreateShareExportPolicyRequestRule{
					{
						Description: sfs.NewNullableString(
							utils.Ptr("first rule"),
						),
						IpAcl:     utils.Ptr([]string{"192.168.2.0/24"}),
						Order:     utils.Ptr(int64(1)),
						SetUuid:   utils.Ptr(true),
						SuperUser: utils.Ptr(false),
					},
					{
						IpAcl:    utils.Ptr([]string{"192.168.2.0/24", "127.0.0.1/32"}),
						Order:    utils.Ptr(int64(2)),
						ReadOnly: utils.Ptr(true),
					},
				}
			}),
		},
	}
	opts := []testutils.TestingOption{
		testutils.WithCmpOptions(cmp.AllowUnexported(sfs.NullableString{})),
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithOptions(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, nil, tt.isValid, opts)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sfs.ApiCreateShareExportPolicyRequest
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
		exportPolicy *sfs.CreateShareExportPolicyResponse
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
			name: "set empty export policy",
			args: args{
				exportPolicy: &sfs.CreateShareExportPolicyResponse{},
			},
			wantErr: true,
		},
		{
			name: "set empty export policy",
			args: args{
				exportPolicy: &sfs.CreateShareExportPolicyResponse{
					ShareExportPolicy: &sfs.CreateShareExportPolicyResponseShareExportPolicy{},
				},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&types.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.projectLabel, tt.args.exportPolicy); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
