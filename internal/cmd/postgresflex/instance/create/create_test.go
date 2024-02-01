package create

import (
	"context"
	"fmt"
	"testing"

	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &postgresflex.APIClient{}

type postgresFlexClientMocked struct {
	listFlavorsFails  bool
	listFlavorsResp   *postgresflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *postgresflex.ListStoragesResponse
}

func (c *postgresFlexClientMocked) CreateInstance(ctx context.Context, projectId string) postgresflex.ApiCreateInstanceRequest {
	return testClient.CreateInstance(ctx, projectId)
}

func (c *postgresFlexClientMocked) ListStoragesExecute(_ context.Context, _, _ string) (*postgresflex.ListStoragesResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return c.listStoragesResp, nil
}

func (c *postgresFlexClientMocked) ListFlavorsExecute(_ context.Context, _ string) (*postgresflex.ListFlavorsResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listFlavorsResp, nil
}

var testProjectId = uuid.NewString()
var testFlavorId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:      testProjectId,
		instanceNameFlag:   "example-name",
		aclFlag:            "0.0.0.0/0",
		backupScheduleFlag: "0 0/6 * * *",
		flavorIdFlag:       testFlavorId,
		storageClassFlag:   "premium-perf4-stackit", // Non-default
		storageSizeFlag:    "10",
		versionFlag:        "6.0",
		typeFlag:           "Replica",
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
		},
		InstanceName:   utils.Ptr("example-name"),
		ACL:            utils.Ptr([]string{"0.0.0.0/0"}),
		BackupSchedule: utils.Ptr("0 0/6 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		StorageClass:   utils.Ptr("premium-perf4-stackit"),
		StorageSize:    utils.Ptr(int64(10)),
		Version:        utils.Ptr("6.0"),
		Type:           utils.Ptr("Replica"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *postgresflex.ApiCreateInstanceRequest)) postgresflex.ApiCreateInstanceRequest {
	request := testClient.CreateInstance(testCtx, testProjectId)
	request = request.CreateInstancePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixturePayload(mods ...func(payload *postgresflex.CreateInstancePayload)) postgresflex.CreateInstancePayload {
	payload := postgresflex.CreateInstancePayload{
		Name:           utils.Ptr("example-name"),
		Acl:            &postgresflex.ACL{Items: utils.Ptr([]string{"0.0.0.0/0"})},
		BackupSchedule: utils.Ptr("0 0/6 * * *"),
		FlavorId:       utils.Ptr(testFlavorId),
		Replicas:       utils.Ptr(int64(3)),
		Storage: &postgresflex.Storage{
			Class: utils.Ptr("premium-perf4-stackit"),
			Size:  utils.Ptr(int64(10)),
		},
		Version: utils.Ptr("6.0"),
		Options: utils.Ptr(map[string]string{
			"type": "Replica",
		}),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
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
				delete(flagValues, versionFlag)
				delete(flagValues, typeFlag)
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
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
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
			cmd := NewCmd()
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			for _, value := range tt.aclValues {
				err := cmd.Flags().Set(aclFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", aclFlag, value, err)
				}
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(cmd, "postgresflex", "create")
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description       string
		model             *inputModel
		expectedRequest   postgresflex.ApiCreateInstanceRequest
		listFlavorsFails  bool
		listFlavorsResp   *postgresflex.ListFlavorsResponse
		listStoragesFails bool
		listStoragesResp  *postgresflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "base with flavor ID",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
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
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
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
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "single instance type",
			model:       fixtureInputModel(func(model *inputModel) { model.Type = utils.Ptr("Single") }),
			isValid:     true,
			expectedRequest: fixtureRequest().CreateInstancePayload(fixturePayload(func(payload *postgresflex.CreateInstancePayload) {
				payload.Options = utils.Ptr(map[string]string{"type": "Single"})
				payload.Replicas = utils.Ptr(int64(1))
			})),
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
		},
		{
			description: "sharded instance type",
			model:       fixtureInputModel(func(model *inputModel) { model.Type = utils.Ptr("Sharded") }),
			isValid:     true,
			expectedRequest: fixtureRequest().CreateInstancePayload(fixturePayload(func(payload *postgresflex.CreateInstancePayload) {
				payload.Options = utils.Ptr(map[string]string{"type": "Sharded"})
				payload.Replicas = utils.Ptr(int64(9))
			})),
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
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
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
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
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
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
			listFlavorsResp: &postgresflex.ListFlavorsResponse{
				Flavors: &[]postgresflex.Flavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"premium-perf4-stackit"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &postgresFlexClientMocked{
				listFlavorsFails:  tt.listFlavorsFails,
				listFlavorsResp:   tt.listFlavorsResp,
				listStoragesFails: tt.listStoragesFails,
				listStoragesResp:  tt.listStoragesResp,
			}
			request, err := buildRequest(testCtx, "postgresflex", tt.model, client)
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
