package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sqlserverflex.APIClient{DefaultAPI: &sqlserverflex.DefaultAPIService{}}

var testRegion = "eu01"

type mockSettings struct {
	listFlavorsFails  bool
	listFlavorsResp   *sqlserverflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *sqlserverflex.ListStoragesResponse
}

func newAPIMock(s mockSettings) sqlserverflex.DefaultAPI {
	return &sqlserverflex.DefaultAPIServiceMock{
		ListStoragesExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListStoragesRequest) (*sqlserverflex.ListStoragesResponse, error) {
			if s.listFlavorsFails {
				return nil, fmt.Errorf("list storages failed")
			}
			return s.listStoragesResp, nil
		}),
		ListFlavorsExecuteMock: utils.Ptr(func(_ sqlserverflex.ApiListFlavorsRequest) (*sqlserverflex.ListFlavorsResponse, error) {
			if s.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return s.listFlavorsResp, nil
		}),
	}
}

var testProjectId = uuid.NewString()
var testFlavorId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		instanceNameFlag:          "example-name",
		aclFlag:                   "0.0.0.0/0",
		backupScheduleFlag:        "0 0/6 * * *",
		flavorIdFlag:              testFlavorId,
		storageClassFlag:          "storage-class", // Non-default
		storageSizeFlag:           "10",
		versionFlag:               "6.0",
		editionFlag:               "developer",
		retentionDaysFlag:         "32",
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
		BackupSchedule: "0 0/6 * * *",
		FlavorId:       utils.Ptr(testFlavorId),
		StorageClass:   "storage-class",
		StorageSize:    utils.Ptr(int64(10)),
		Version:        "6.0",
		RetentionDays:  utils.Ptr(int32(32)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sqlserverflex.ApiCreateInstanceRequest)) sqlserverflex.ApiCreateInstanceRequest {
	request := testClient.DefaultAPI.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *sqlserverflex.CreateInstancePayload)) sqlserverflex.CreateInstancePayload {
	payload := sqlserverflex.CreateInstancePayload{
		Name: "example-name",
		Network: sqlserverflex.CreateInstancePayloadNetwork{
			Acl: []string{"0.0.0.0/0"},
		},
		BackupSchedule: "0 0/6 * * *",
		FlavorId:       testFlavorId,
		Storage: sqlserverflex.StorageCreate{
			Class: "storage-class",
			Size:  int64(10),
		},
		Version:       "6.0",
		RetentionDays: 32,
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
			isValid: false,
		},
		{
			description: "repeated acl flags",
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL = append(model.ACL, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "repeated acl flag with list value",
			flagValues:  fixtureFlagValues(),
			aclValues:   []string{"198.51.100.14/24,198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL = append(model.ACL, "198.51.100.14/24", "198.51.100.14/32")
			}),
		},
		{
			description: "no acls",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			aclValues: []string{},
			isValid:   true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.ACL = nil
			}),
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
		expectedRequest   sqlserverflex.ApiCreateInstanceRequest
		listFlavorsFails  bool
		listFlavorsResp   *sqlserverflex.ListFlavorsResponse
		listStoragesFails bool
		listStoragesResp  *sqlserverflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "base with flavor ID",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id:     testFlavorId,
						Cpu:    int64(2),
						Memory: int64(4),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{{
					Class: "storage-class",
				}},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: 10,
					Max: 100,
				},
			},
		},
		{
			description: "with CPU and RAM",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id:     testFlavorId,
						Cpu:    int64(2),
						Memory: int64(4),
					},
					{
						Id:     "other-flavor",
						Cpu:    int64(1),
						Memory: int64(8),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{{
					Class: "storage-class",
				}},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(10),
					Max: int32(100),
				},
			},
		},
		{
			description: "get flavors fails",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
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
					model.CPU = utils.Ptr(int64(5))
					model.RAM = utils.Ptr(int64(9))
				},
			),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id:     testFlavorId,
						Cpu:    int64(2),
						Memory: int64(4),
					},
					{
						Id:     "other-flavor",
						Cpu:    int64(1),
						Memory: int64(8),
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
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "invalid storage class",
			model: fixtureInputModel(
				func(model *inputModel) {
					model.StorageClass = "non-existing-class"
				},
			),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id:     testFlavorId,
						Cpu:    int64(2),
						Memory: int64(4),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{{
					Class: "storage-class",
				}},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(10),
					Max: int32(100),
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
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: []sqlserverflex.ListFlavors{
					{
						Id:     testFlavorId,
						Cpu:    int64(2),
						Memory: int64(4),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: []sqlserverflex.FlavorStorageClassesStorageClass{{
					Class: "storage-class",
				}},
				StorageRange: sqlserverflex.FlavorStorageRange{
					Min: int32(10),
					Max: int32(100),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := mockSettings{
				listFlavorsFails:  tt.listFlavorsFails,
				listFlavorsResp:   tt.listFlavorsResp,
				listStoragesFails: tt.listStoragesFails,
				listStoragesResp:  tt.listStoragesResp,
			}
			request, err := buildRequest(testCtx, tt.model, newAPIMock(client))
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
				cmp.FilterPath(func(p cmp.Path) bool {
					return p.String() == "ApiService"
				}, cmp.Ignore()),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		model        *inputModel
		projectLabel string
		resp         *sqlserverflex.CreateInstanceResponse
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
			name: "sql instance as argument",
			args: args{
				model: fixtureInputModel(),
				resp:  &sqlserverflex.CreateInstanceResponse{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
