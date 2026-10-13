package clone

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &postgresflex.APIClient{DefaultAPI: &postgresflex.DefaultAPIService{}}

	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
)

const (
	testStorageSize       = int64(10)
	testRecoveryTimestamp = "2024-03-08T09:28:00+00:00"
	testStorageClass      = "premium-perf4-stackit"
	testRegion            = "eu01"
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

func fixtureRequiredFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		recoveryTimestampFlag:     testRecoveryTimestamp,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureStandardFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		recoveryTimestampFlag:     testRecoveryTimestamp,
		storageClassFlag:          "class",
		storageSizeFlag:           "10",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureRequiredInputModel(mods ...func(model *inputModel)) *inputModel {
	testRecoveryTimestamp, err := time.Parse(recoveryDateFormat, testRecoveryTimestamp)
	if err != nil {
		return &inputModel{}
	}

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:   testInstanceId,
		RecoveryDate: testRecoveryTimestamp,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureStandardInputModel(mods ...func(model *inputModel)) *inputModel {
	testRecoveryTimestamp, err := time.Parse(recoveryDateFormat, testRecoveryTimestamp)
	if err != nil {
		return &inputModel{}
	}

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:   testInstanceId,
		StorageClass: utils.Ptr(testStorageClass),
		StorageSize:  utils.Ptr(testStorageSize),
		RecoveryDate: testRecoveryTimestamp,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *postgresflex.ApiCloneInstanceRequest)) postgresflex.ApiCloneInstanceRequest {
	request := testClient.DefaultAPI.CloneInstance(testCtx, testProjectId, testRegion, testInstanceId)
	request = request.CloneInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *postgresflex.CloneInstancePayload)) postgresflex.CloneInstancePayload {
	testRecoveryTimestamp, err := time.Parse(recoveryDateFormat, testRecoveryTimestamp)
	if err != nil {
		return postgresflex.CloneInstancePayload{}
	}

	payload := postgresflex.CloneInstancePayload{
		InstanceOverrides: postgresflex.CloneInstanceOverrides{
			Class: testStorageClass,
			Size:  testStorageSize,
		},
		PointInTime: testRecoveryTimestamp,
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
			flagValues:    fixtureRequiredFlagValues(),
			isValid:       true,
			expectedModel: fixtureRequiredInputModel(),
		},
		{
			description: "with defaults",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, storageClassFlag)
				delete(flagValues, storageSizeFlag)
			}),
			isValid:       true,
			expectedModel: fixtureRequiredInputModel(),
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
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "all values with storage class",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, storageSizeFlag)
				flagValues[storageClassFlag] = "premium-perf4-stackit"
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.StorageSize = nil
				model.StorageClass = utils.Ptr("premium-perf4-stackit")
			}),
		},
		{
			description: "all values with storage size",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, storageClassFlag)
				flagValues[storageSizeFlag] = "2"
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.StorageClass = nil
				model.StorageSize = utils.Ptr(int64(2))
			}),
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "recovery timestamp is missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				delete(flagValues, recoveryTimestampFlag)
			}),
			isValid: false,
		},
		{
			description: "recovery timestamp is empty",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[recoveryTimestampFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "recovery timestamp is invalid",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[recoveryTimestampFlag] = "test"
			}),
			isValid: false,
		},
		{
			description: "recovery timestamp is invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[recoveryTimestampFlag] = "11:00 12/12/2024"
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
	testRecoveryTimestamp, err := time.Parse(recoveryDateFormat, testRecoveryTimestamp)
	if err != nil {
		return
	}

	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest postgresflex.ApiCloneInstanceRequest
		isValid         bool
	}{
		{
			description: "base",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr(testStorageClass)
					model.StorageSize = utils.Ptr(testStorageSize)
				},
			),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "specify storage class only",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageSize = utils.Ptr(testStorageSize)
				model.StorageClass = utils.Ptr("class")
			}),
			isValid: true,
			expectedRequest: testClient.DefaultAPI.CloneInstance(testCtx, testProjectId, testRegion, testInstanceId).
				CloneInstancePayload(postgresflex.CloneInstancePayload{
					InstanceOverrides: postgresflex.CloneInstanceOverrides{
						Class: "class",
						Size:  testStorageSize,
					},
					PointInTime: testRecoveryTimestamp,
				}),
		},
		{
			description: "storage size missing",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr(testStorageClass)
				model.StorageSize = nil
			}),
			isValid: false,
		},
		{
			description: "storage class missing",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = nil
				model.StorageSize = utils.Ptr(testStorageSize)
			}),
			isValid: false,
		},
		{
			description: "specify storage class and size",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr("class")
				model.StorageSize = utils.Ptr(int64(10))
			}),
			isValid: true,
			expectedRequest: testClient.DefaultAPI.CloneInstance(testCtx, testProjectId, testRegion, testInstanceId).
				CloneInstancePayload(postgresflex.CloneInstancePayload{
					InstanceOverrides: postgresflex.CloneInstanceOverrides{
						Class: "class",
						Size:  int64(10),
					},
					PointInTime: testRecoveryTimestamp,
				}),
		},
		{
			description: "get instance fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("class")
					model.RecoveryDate = testRecoveryTimestamp
				},
			),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient.DefaultAPI)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmpopts.IgnoreFields(tt.expectedRequest, "ApiService"),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		OutputFormat  string
		instanceLabel string
		instanceId    string
		async         bool
		resp          *postgresflex.CloneInstanceResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"standard", args{
			instanceLabel: "foo",
			instanceId:    "bar",
			resp:          &postgresflex.CloneInstanceResponse{Id: "id"},
		}, false},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.OutputFormat, tt.args.async, tt.args.instanceLabel, tt.args.instanceId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
