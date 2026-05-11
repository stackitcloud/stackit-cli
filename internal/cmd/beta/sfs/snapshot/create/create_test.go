package create

import (
	"context"
	"testing"
	"time"

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
var testComment = "test-comment"
var testResourcePoolId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,

		nameFlag:           testName,
		resourcePoolIdFlag: testResourcePoolId,
		commentFlag:        testComment,
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
		Name:           testName,
		ResourcePoolId: testResourcePoolId,
		Comment:        utils.Ptr(testComment),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiCreateResourcePoolSnapshotRequest)) sfs.ApiCreateResourcePoolSnapshotRequest {
	request := testClient.DefaultAPI.CreateResourcePoolSnapshot(testCtx, testProjectId, testRegion, testResourcePoolId)
	request = request.CreateResourcePoolSnapshotPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(request *sfs.CreateResourcePoolSnapshotPayload)) sfs.CreateResourcePoolSnapshotPayload {
	payload := sfs.CreateResourcePoolSnapshotPayload{
		Name: utils.Ptr(testName),
		Comment: *sfs.NewNullableString(
			utils.Ptr(testComment),
		),
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
				delete(flagValues, commentFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Comment = nil
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
		{
			description: "invalid resource pool id 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[resourcePoolIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid resource pool id 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[resourcePoolIdFlag] = "invalid-resource-pool-id"
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
		expectedRequest sfs.ApiCreateResourcePoolSnapshotRequest
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
				cmp.AllowUnexported(sfs.NullableString{}, sfs.NullableInt32{}),
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
		snapshotLabel     string
		resourcePoolLabel string
		resp              *sfs.CreateResourcePoolSnapshotResponse
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
				resp: &sfs.CreateResourcePoolSnapshotResponse{},
			},
			wantErr: false,
		},
		{
			name: "set empty snapshot",
			args: args{
				resp: &sfs.CreateResourcePoolSnapshotResponse{
					ResourcePoolSnapshot: &sfs.ResourcePoolSnapshot{},
				},
			},
			wantErr: false,
		},
		{
			name: "set full snapshot",
			args: args{
				resp: &sfs.CreateResourcePoolSnapshotResponse{
					ResourcePoolSnapshot: &sfs.ResourcePoolSnapshot{
						SnaplockExpiryTime: *sfs.NewNullableTime(utils.Ptr(time.Now().Add(time.Hour))),
					},
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.snapshotLabel, tt.args.resourcePoolLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
