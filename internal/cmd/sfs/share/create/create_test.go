package create

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sfs "github.com/stackitcloud/stackit-sdk-go/services/sfs/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var projectIdFlag = globalflags.ProjectIdFlag
var regionFlag = globalflags.RegionFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sfs.APIClient{DefaultAPI: &sfs.DefaultAPIService{}}

var testProjectId = uuid.NewString()
var testRegion = "eu01"

var testName = "test-name"
var testResourcePoolId = uuid.NewString()
var testExportPolicyName = "test-export-policy"
var testHardLimit int32 = 10

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,

		nameFlag:             testName,
		resourcePoolIdFlag:   testResourcePoolId,
		exportPolicyNameFlag: testExportPolicyName,
		hardLimitFlag:        strconv.Itoa(int(testHardLimit)),
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
		Name:             testName,
		ResourcePoolId:   testResourcePoolId,
		ExportPolicyName: utils.Ptr(testExportPolicyName),
		HardLimit:        utils.Ptr(testHardLimit),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiCreateShareRequest)) sfs.ApiCreateShareRequest {
	request := testClient.DefaultAPI.CreateShare(testCtx, testProjectId, testRegion, testResourcePoolId)
	request = request.CreateSharePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(request *sfs.CreateSharePayload)) sfs.CreateSharePayload {
	payload := sfs.CreateSharePayload{
		Name:                    testName,
		ExportPolicyName:        *sfs.NewNullableString(utils.Ptr(testExportPolicyName)),
		SpaceHardLimitGigabytes: testHardLimit,
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
				delete(flagValues, exportPolicyNameFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ExportPolicyName = nil
			}),
		},
		{
			description: "missing required name",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "missing required resourcePoolId",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, resourcePoolIdFlag)
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
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest sfs.ApiCreateShareRequest
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
				cmp.AllowUnexported(tt.expectedRequest, sfs.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(sfs.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat      string
		async             bool
		resourcePoolLabel string
		item              *sfs.CreateShareResponse
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
			name: "set empty response",
			args: args{
				item: &sfs.CreateShareResponse{},
			},
			wantErr: false,
		},
		{
			name: "set empty response share",
			args: args{
				item: &sfs.CreateShareResponse{
					Share: &sfs.Share{},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.resourcePoolLabel, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
