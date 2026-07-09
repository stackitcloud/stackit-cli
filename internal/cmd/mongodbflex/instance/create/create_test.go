package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{DefaultAPI: &mongodbflex.DefaultAPIService{}}

type mockSettings struct {
	listFlavorsFails  bool
	listFlavorsResp   *mongodbflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *mongodbflex.ListStoragesResponse
}

var testProjectId = uuid.NewString()
var testFlavorId = uuid.NewString()

func newAPICLientMock(settings mockSettings) mongodbflex.DefaultAPI {
	return mongodbflex.DefaultAPIServiceMock{
		ListStoragesExecuteMock: utils.Ptr(func(_ mongodbflex.ApiListStoragesRequest) (*mongodbflex.ListStoragesResponse, error) {
			if settings.listFlavorsFails {
				return nil, fmt.Errorf("list storages failed")
			}
			return settings.listStoragesResp, nil
		}),
		ListFlavorsExecuteMock: utils.Ptr(func(_ mongodbflex.ApiListFlavorsRequest) (*mongodbflex.ListFlavorsResponse, error) {
			if settings.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return settings.listFlavorsResp, nil
		}),
	}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceNameFlag:          "example-name",
		aclFlag:                   "0.0.0.0/0",
		backupScheduleFlag:        "0 0/6 * * *",
		flavorIdFlag:              testFlavorId,
		storageClassFlag:          "premium-perf4-mongodb", // Non-default
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
		InstanceName:   utils.Ptr("example-name"),
		ACL:            utils.Ptr([]string{"0.0.0.0/0"}),
		BackupSchedule: utils.Ptr("0 0/6 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		StorageClass:   utils.Ptr("premium-perf4-mongodb"),
		StorageSize:    utils.Ptr(int64(10)),
		Version:        utils.Ptr("6.0"),
		Type:           utils.Ptr("Replica"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *mongodbflex.ApiCreateInstanceRequest)) mongodbflex.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *mongodbflex.CreateInstancePayload)) mongodbflex.CreateInstancePayload {
	payload := mongodbflex.CreateInstancePayload{
		Name:           "example-name",
		Acl:            mongodbflex.ACL{Items: []string{"0.0.0.0/0"}},
		BackupSchedule: "0 0/6 * * *",
		FlavorId:       testFlavorId,
		Replicas:       int32(3),
		Storage: mongodbflex.Storage{
			Class: utils.Ptr("premium-perf4-mongodb"),
			Size:  utils.Ptr(int64(10)),
		},
		Version: "6.0",
		Options: map[string]string{
			"type": "Replica",
		},
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
			isValid:       true,
			expectedModel: fixtureInputModel(),
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
				model.CPU = utils.Ptr(int32(2))
				model.RAM = utils.Ptr(int32(4))
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
				model.ACL = utils.Ptr(
					append(*model.ACL, "198.51.100.14/24", "198.51.100.14/32"),
				)
			}),
		},
		{
			description: "repeated acl flag with list value",
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL = utils.Ptr(
					append(*model.ACL, "198.51.100.14/24", "198.51.100.14/32"),
				)
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
		description       string
		model             *inputModel
		expectedRequest   mongodbflex.ApiCreateInstanceRequest
		listFlavorsFails  bool
		listFlavorsResp   *mongodbflex.ListFlavorsResponse
		listStoragesFails bool
		listStoragesResp  *mongodbflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "base with flavor ID",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "with CPU and RAM",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int32(2))
					model.RAM = utils.Ptr(int32(4))
				},
			),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
					{
						Id:     utils.Ptr("other-flavor"),
						Cpu:    utils.Ptr(int32(1)),
						Memory: utils.Ptr(int32(8)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "single instance type",
			model:       fixtureInputModel(func(model *inputModel) { model.Type = utils.Ptr("Single") }),
			isValid:     true,
			expectedRequest: fixtureRequest().CreateInstancePayload(fixturePayload(func(payload *mongodbflex.CreateInstancePayload) {
				payload.Options = map[string]string{"type": "Single"}
				payload.Replicas = int32(1)
			})),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "sharded instance type",
			model:       fixtureInputModel(func(model *inputModel) { model.Type = utils.Ptr("Sharded") }),
			isValid:     true,
			expectedRequest: fixtureRequest().CreateInstancePayload(fixturePayload(func(payload *mongodbflex.CreateInstancePayload) {
				payload.Options = map[string]string{"type": "Sharded"}
				payload.Replicas = int32(9)
			})),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "get flavors fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int32(2))
					model.RAM = utils.Ptr(int32(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "flavor id not found",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int32(5))
					model.RAM = utils.Ptr(int32(9))
				},
			),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
					{
						Id:     utils.Ptr("other-flavor"),
						Cpu:    utils.Ptr(int32(1)),
						Memory: utils.Ptr(int32(8)),
					},
				},
			},
			isValid: false,
		},
		{
			description: "get storages fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int32(2))
					model.RAM = utils.Ptr(int32(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "invalid storage class",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("non-existing-class")
				},
			),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
		{
			description: "invalid storage size",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.StorageSize = utils.Ptr(int64(9))
				},
			),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: []mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int32(2)),
						Memory: utils.Ptr(int32(4)),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: []string{"premium-perf4-mongodb"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			settings := mockSettings{
				listFlavorsFails:  tt.listFlavorsFails,
				listFlavorsResp:   tt.listFlavorsResp,
				listStoragesFails: tt.listStoragesFails,
				listStoragesResp:  tt.listStoragesResp,
			}
			request, err := buildRequest(testCtx, tt.model, newAPICLientMock(settings))
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, mongodbflex.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
				cmpopts.IgnoreFields(tt.expectedRequest, "ApiService"),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat           string
		async                  bool
		projectLabel           string
		createInstanceResponse *mongodbflex.CreateInstanceResponse
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
			name: "set empty create instance response",
			args: args{
				createInstanceResponse: &mongodbflex.CreateInstanceResponse{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.projectLabel, tt.args.createInstanceResponse); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
