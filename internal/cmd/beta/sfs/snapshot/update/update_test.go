package update

import (
	"context"
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

var testSnapshotName = "test-snapshot-name"
var testNewName = "test-new-name"
var testComment = "test-comment"
var testResourcePoolId = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testSnapshotName,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag: testProjectId,
		regionFlag:    testRegion,

		newSnapshotNameFlag: testNewName,
		resourcePoolIdFlag:  testResourcePoolId,
		commentFlag:         testComment,
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
		SnapshotName:    testSnapshotName,
		NewSnapshotName: testNewName,
		ResourcePoolId:  testResourcePoolId,
		Comment:         utils.Ptr(testComment),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sfs.ApiUpdateResourcePoolSnapshotRequest)) sfs.ApiUpdateResourcePoolSnapshotRequest {
	request := testClient.DefaultAPI.UpdateResourcePoolSnapshot(testCtx, testProjectId, testRegion, testResourcePoolId, testSnapshotName)
	request = request.UpdateResourcePoolSnapshotPayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(request *sfs.UpdateResourcePoolSnapshotPayload)) sfs.UpdateResourcePoolSnapshotPayload {
	payload := sfs.UpdateResourcePoolSnapshotPayload{
		Name: *sfs.NewNullableString(utils.Ptr(testNewName)),
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "either name or comment (only name set)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, commentFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Comment = nil
			}),
		},
		{
			description: "either name or comment (only comment set)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, newSnapshotNameFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.NewSnapshotName = ""
			}),
		},
		{
			description: "missing both name and comment (at least one required)",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, newSnapshotNameFlag)
				delete(flagValues, commentFlag)
			}),
			isValid: false,
		},
		{
			description: "missing required resourcePoolId",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, resourcePoolIdFlag)
			}),
			isValid: false,
		},
		{
			description: "invalid resource pool id 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[resourcePoolIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid resource pool id 2",
			argValues:   fixtureArgValues(),
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
		expectedRequest sfs.ApiUpdateResourcePoolSnapshotRequest
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
		resp              *sfs.UpdateResourcePoolSnapshotResponse
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
				resp: &sfs.UpdateResourcePoolSnapshotResponse{},
			},
			wantErr: false,
		},
		{
			name: "set empty snapshot",
			args: args{
				resp: &sfs.UpdateResourcePoolSnapshotResponse{
					ResourcePoolSnapshot: &sfs.ResourcePoolSnapshot{},
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
