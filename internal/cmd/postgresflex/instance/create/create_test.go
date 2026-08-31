package create

import (
	"context"
	"testing"

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

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &postgresflex.APIClient{DefaultAPI: &postgresflex.DefaultAPIService{}}
var testProjectId = uuid.NewString()
var testFlavorId = uuid.NewString()

const testRegion = "eu01"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceNameFlag:          "example-name",
		aclFlag:                   "0.0.0.0/0",
		backupScheduleFlag:        "0 0 * * *",
		flavorIdFlag:              testFlavorId,
		storageClassFlag:          "premium-perf4-stackit", // Non-default
		storageSizeFlag:           "10",
		versionFlag:               "6.0",
		typeFlag.Name():           "Replica",
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
		InstanceName:   "example-name",
		ACL:            []string{"0.0.0.0/0"},
		BackupSchedule: utils.Ptr("0 0 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		StorageClass:   utils.Ptr("premium-perf4-stackit"),
		StorageSize:    utils.Ptr(int64(10)),
		Version:        utils.Ptr("6.0"),
		Type:           "Replica",
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *postgresflex.ApiCreateInstanceRequest)) postgresflex.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *postgresflex.CreateInstancePayload)) postgresflex.CreateInstancePayload {
	payload := postgresflex.CreateInstancePayload{
		BackupSchedule: "0 0 * * *",
		FlavorId:       testFlavorId,
		Name:           "example-name",
		Network: postgresflex.InstanceNetworkCreate{
			Acl: []string{"0.0.0.0/0"},
		},
		Storage: postgresflex.StorageCreate{
			Class: "premium-perf4-stackit",
			Size:  int64(10),
		},
		Version:       "6.0",
		RetentionDays: *postgresflex.NewNullableInt32(nil),
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
		aclValues     []string
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
			description: "with defaults",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, backupScheduleFlag)
				delete(flagValues, typeFlag.Name())
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.BackupSchedule = nil
			}),
		},
		{
			description: "use CPU and RAM",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
				delete(flagValues, flavorIdFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.FlavorId = nil
				model.CPU = utils.Ptr(int64(2))
				model.RAM = utils.Ptr(int64(4))
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
			description: "invalid with flavor ID, CPU and RAM",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
			}),
			isValid: false,
		},
		{
			description: "invalid with flavor ID and CPU",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "2"
			}),
			isValid: false,
		},
		{
			description: "invalid with CPU only",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[cpuFlag] = "2"
			}),
			isValid: false,
		},
		{
			description: "no version",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, versionFlag)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Version = nil
			}),
		},
		{
			description: "repeated acl flags",
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL =
					append(model.ACL, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "repeated acl flag with list value",
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL =
					append(model.ACL, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "no acls",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			aclValues: []string{},
			isValid:   false,
		},
		{
			description: "access scope set",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues["access-scope"] = string(postgresflex.INSTANCENETWORKACCESSSCOPE_SNA)
			}),
			isValid: true,
			expectedModel: fixtureInputModel(
				func(model *inputModel) {
					model.AccessScope = new(string(postgresflex.INSTANCENETWORKACCESSSCOPE_SNA))
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInputWithAdditionalFlags(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, map[string][]string{
				aclFlag: tt.aclValues,
			}, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest postgresflex.ApiCreateInstanceRequest
		isValid         bool
	}{
		{
			description:     "base with flavor ID",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "access scope set",
			model: fixtureInputModel(func(model *inputModel) {
				model.AccessScope = new(string(postgresflex.INSTANCENETWORKACCESSSCOPE_SNA))
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *postgresflex.ApiCreateInstanceRequest) {
				payload := fixturePayload(func(payload *postgresflex.CreateInstancePayload) {
					payload.Network.AccessScope = new(postgresflex.INSTANCENETWORKACCESSSCOPE_SNA)
				})

				*request = request.CreateInstancePayload(payload)
			}),
		},
		{
			description: "storage class missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.StorageClass = nil
			}),
			isValid: false,
		},
		{
			description: "storage size missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.StorageSize = nil
			}),
			isValid: false,
		},
		{
			description: "backup schedule missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.BackupSchedule = nil
			}),
			isValid: false,
		},
		{
			description: "version missing",
			model: fixtureInputModel(func(model *inputModel) {
				model.Version = nil
			}),
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
				cmpopts.EquateComparable(testCtx, postgresflex.DefaultAPIService{}, postgresflex.NullableInt32{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func Test_outputResult(t *testing.T) {
	type args struct {
		outputFormat string
		async        bool
		projectLabel string
		instanceId   string
		resp         *postgresflex.CreateInstanceResponse
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
			name: "standard",
			args: args{
				outputFormat: "",
				async:        false,
				projectLabel: "label",
				instanceId:   "4711",
				resp: &postgresflex.CreateInstanceResponse{
					Id: "id",
				},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.projectLabel, tt.args.instanceId, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
