package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sqlserverflex.APIClient{}
var testRegion = "eu01"

// enforce implementation of interfaces
var (
	_ sqlServerFlexClient = &sqlServerFlexClientMocked{}
)

type sqlServerFlexClientMocked struct {
	listFlavorsFails  bool
	listFlavorsResp   *sqlserverflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *sqlserverflex.ListStoragesResponse
}

func (c *sqlServerFlexClientMocked) CreateInstance(ctx context.Context, projectId, region string) sqlserverflex.ApiCreateInstanceRequest {
	return testClient.CreateInstance(ctx, projectId, region)
}

func (c *sqlServerFlexClientMocked) ListStoragesExecute(_ context.Context, _, _, _ string) (*sqlserverflex.ListStoragesResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return c.listStoragesResp, nil
}

func (c *sqlServerFlexClientMocked) ListFlavorsExecute(_ context.Context, _, _ string) (*sqlserverflex.ListFlavorsResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listFlavorsResp, nil
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
		InstanceName:   utils.Ptr("example-name"),
		ACL:            utils.Ptr([]string{"0.0.0.0/0"}),
		BackupSchedule: utils.Ptr("0 0/6 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		StorageClass:   utils.Ptr("storage-class"),
		StorageSize:    utils.Ptr(int64(10)),
		Version:        utils.Ptr("6.0"),
		Edition:        utils.Ptr("developer"),
		RetentionDays:  utils.Ptr(int64(32)),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sqlserverflex.ApiCreateInstanceRequest)) sqlserverflex.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId, testRegion)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *sqlserverflex.CreateInstancePayload)) sqlserverflex.CreateInstancePayload {
	payload := sqlserverflex.CreateInstancePayload{
		Name:           utils.Ptr("example-name"),
		Acl:            &sqlserverflex.CreateInstancePayloadAcl{Items: utils.Ptr([]string{"0.0.0.0/0"})},
		BackupSchedule: utils.Ptr("0 0/6 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		Storage: &sqlserverflex.CreateInstancePayloadStorage{
			Class: utils.Ptr("storage-class"),
			Size:  utils.Ptr(int64(10)),
		},
		Version: utils.Ptr("6.0"),
		Options: &sqlserverflex.CreateInstancePayloadOptions{
			Edition:       utils.Ptr("developer"),
			RetentionDays: utils.Ptr("32"),
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
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: &[]string{"storage-class"},
				StorageRange: &sqlserverflex.StorageRange{
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
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
					{
						Id:     utils.Ptr("other-flavor"),
						Cpu:    utils.Ptr(int64(1)),
						Memory: utils.Ptr(int64(8)),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: &[]string{"storage-class"},
				StorageRange: &sqlserverflex.StorageRange{
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
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
					{
						Id:     utils.Ptr("other-flavor"),
						Cpu:    utils.Ptr(int64(1)),
						Memory: utils.Ptr(int64(8)),
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
					model.StorageClass = utils.Ptr("non-existing-class")
				},
			),
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: &[]string{"storage-class"},
				StorageRange: &sqlserverflex.StorageRange{
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
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &sqlserverflex.ListStoragesResponse{
				StorageClasses: &[]string{"storage-class"},
				StorageRange: &sqlserverflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &sqlServerFlexClientMocked{
				listFlavorsFails:  tt.listFlavorsFails,
				listFlavorsResp:   tt.listFlavorsResp,
				listStoragesFails: tt.listStoragesFails,
				listStoragesResp:  tt.listStoragesResp,
			}
			request, err := buildRequest(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

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
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.model, tt.args.projectLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
