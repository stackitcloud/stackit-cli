package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &sqlserverflex.APIClient{}

type mongoDBFlexClientMocked struct {
	listFlavorsFails  bool
	listFlavorsResp   *sqlserverflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *sqlserverflex.ListStoragesResponse
	getInstanceFails  bool
	getInstanceResp   *sqlserverflex.GetInstanceResponse
}

func (c *mongoDBFlexClientMocked) PartialUpdateInstance(ctx context.Context, projectId, instanceId string) sqlserverflex.ApiPartialUpdateInstanceRequest {
	return testClient.PartialUpdateInstance(ctx, projectId, instanceId)
}

func (c *mongoDBFlexClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*sqlserverflex.GetInstanceResponse, error) {
	if c.getInstanceFails {
		return nil, fmt.Errorf("get instance failed")
	}
	return c.getInstanceResp, nil
}

func (c *mongoDBFlexClientMocked) ListStoragesExecute(_ context.Context, _, _ string) (*sqlserverflex.ListStoragesResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return c.listStoragesResp, nil
}

func (c *mongoDBFlexClientMocked) ListFlavorsExecute(_ context.Context, _ string) (*sqlserverflex.ListFlavorsResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listFlavorsResp, nil
}

var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testFlavorId = uuid.NewString()

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
		projectIdFlag: testProjectId,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureStandardFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:      testProjectId,
		flavorIdFlag:       testFlavorId,
		instanceNameFlag:   "example-name",
		aclFlag:            "0.0.0.0/0",
		backupScheduleFlag: "0 0 * * *",
		versionFlag:        "5.0",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureRequiredInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId: testInstanceId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureStandardInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:     testInstanceId,
		FlavorId:       utils.Ptr(testFlavorId),
		InstanceName:   utils.Ptr("example-name"),
		ACL:            utils.Ptr([]string{"0.0.0.0/0"}),
		BackupSchedule: utils.Ptr("0 0 * * *"),
		Version:        utils.Ptr("5.0"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *sqlserverflex.ApiPartialUpdateInstanceRequest)) sqlserverflex.ApiPartialUpdateInstanceRequest {
	request := testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId)
	request = request.PartialUpdateInstancePayload(sqlserverflex.PartialUpdateInstancePayload{})
	for _, mod := range mods {
		mod(&request)
	}
	return request
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
			description: "only instance and project ids",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureRequiredFlagValues(),

			isValid: false,
		},
		{
			description:   "all values with flavor id",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureStandardFlagValues(),
			isValid:       true,
			expectedModel: fixtureStandardInputModel(),
		},
		{
			description: "all values with cpu and ram",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, flavorIdFlag)
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.FlavorId = nil
				model.CPU = utils.Ptr(int64(2))
				model.RAM = utils.Ptr(int64(4))
			}),
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				delete(flagValues, projectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[projectIdFlag] = "invalid-uuid"
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
			description: "invalid with flavor ID, CPU and RAM",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[flavorIdFlag] = testFlavorId
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
			}),
			isValid: false,
		},
		{
			description: "invalid with flavor ID and CPU",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[flavorIdFlag] = testFlavorId
				flagValues[cpuFlag] = "2"
			}),
			isValid: false,
		},
		{
			description: "no acl flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.ACL = nil
			}),
		},
		{
			description: "repeated acl flags",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureRequiredFlagValues(),
			aclValues:   []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureRequiredInputModel(func(model *inputModel) {
				model.ACL = utils.Ptr([]string{"198.51.100.14/24", "198.51.100.14/32"})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
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

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
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
		expectedRequest   sqlserverflex.ApiPartialUpdateInstanceRequest
		getInstanceFails  bool
		getInstanceResp   *sqlserverflex.GetInstanceResponse
		listFlavorsFails  bool
		listFlavorsResp   *sqlserverflex.ListFlavorsResponse
		listStoragesFails bool
		listStoragesResp  *sqlserverflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "no values",
			model:           fixtureRequiredInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "update flavor from id",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.FlavorId = utils.Ptr(testFlavorId)
			}),
			isValid: true,
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId).
				PartialUpdateInstancePayload(sqlserverflex.PartialUpdateInstancePayload{
					FlavorId: utils.Ptr(testFlavorId),
				}),
		},
		{
			description: "update flavor from cpu and ram",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.CPU = utils.Ptr(int64(2))
				model.RAM = utils.Ptr(int64(4))
			}),
			isValid: true,
			listFlavorsResp: &sqlserverflex.ListFlavorsResponse{
				Flavors: &[]sqlserverflex.InstanceFlavorEntry{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId).
				PartialUpdateInstancePayload(sqlserverflex.PartialUpdateInstancePayload{
					FlavorId: utils.Ptr(testFlavorId),
				}),
		},
		{
			description: "get flavors fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "flavor id not found",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
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
			description: "get instance fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.RAM = utils.Ptr(int64(4))
				},
			),
			getInstanceFails: true,
			isValid:          false,
		},
		{
			description: "get storages fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &mongoDBFlexClientMocked{
				getInstanceFails:  tt.getInstanceFails,
				getInstanceResp:   tt.getInstanceResp,
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
			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
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
