package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
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

var (
	testCtx    = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient = &mongodbflex.APIClient{DefaultAPI: &mongodbflex.DefaultAPIService{}}

	testProjectId  = uuid.NewString()
	testInstanceId = uuid.NewString()
	testFlavorId   = uuid.NewString()
)

type mockClientSettings struct {
	listFlavorsFails bool
	listFlavorsResp  *mongodbflex.ListFlavorsResponse
	getInstanceFails bool
	getInstanceResp  *mongodbflex.InstanceResponse
}

func newAPIClientMock(c mockClientSettings) mongodbflex.DefaultAPI {
	return mongodbflex.DefaultAPIServiceMock{
		GetInstanceExecuteMock: utils.Ptr(func(_ mongodbflex.ApiGetInstanceRequest) (*mongodbflex.InstanceResponse, error) {
			if c.getInstanceFails {
				return nil, fmt.Errorf("get instance failed")
			}
			return c.getInstanceResp, nil
		}),
		ListFlavorsExecuteMock: utils.Ptr(func(_ mongodbflex.ApiListFlavorsRequest) (*mongodbflex.ListFlavorsResponse, error) {
			if c.listFlavorsFails {
				return nil, fmt.Errorf("list flavors failed")
			}
			return c.listFlavorsResp, nil
		}),
	}
}

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
		flavorIdFlag:              testFlavorId,
		instanceNameFlag:          "example-name",
		aclFlag:                   "0.0.0.0/0",
		backupScheduleFlag:        "0 0 * * *",
		storageClassFlag:          "class",
		storageSizeFlag:           "10",
		versionFlag:               "5.0",
		typeFlag.Name():           "Single",
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
			Region:    testRegion,
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
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:     testInstanceId,
		FlavorId:       new(testFlavorId),
		InstanceName:   new("example-name"),
		ACL:            new([]string{"0.0.0.0/0"}),
		BackupSchedule: new("0 0 * * *"),
		StorageClass:   new("class"),
		StorageSize:    new(int64(10)),
		Version:        new("5.0"),
		Type:           new("Single"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *mongodbflex.ApiPartialUpdateInstanceRequest)) mongodbflex.ApiPartialUpdateInstanceRequest {
	request := testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion)
	request = request.PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{})
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
				model.CPU = new(int32(2))
				model.RAM = new(int32(4))
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
				model.ACL = new([]string{"198.51.100.14/24", "198.51.100.14/32"})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			typeFlag.Reset()
			params := testparams.NewTestParams()
			cmd := NewCmd(params.CmdParams)
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

			model, err := parseInput(params.Printer, cmd, tt.argValues)
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
		description        string
		model              *inputModel
		expectedRequest    mongodbflex.ApiPartialUpdateInstanceRequest
		mockClientSettings mockClientSettings
		isValid            bool
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
				model.FlavorId = new(testFlavorId)
			}),
			isValid: true,
			mockClientSettings: mockClientSettings{
				listFlavorsResp: &mongodbflex.ListFlavorsResponse{
					Flavors: []mongodbflex.InstanceFlavor{
						{
							Id:     new(testFlavorId),
							Cpu:    new(int32(2)),
							Memory: new(int32(4)),
						},
					},
				},
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					FlavorId: new(testFlavorId),
				}),
		},
		{
			description: "update flavor from cpu and ram",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.CPU = new(int32(2))
				model.RAM = new(int32(4))
			}),
			isValid: true,
			mockClientSettings: mockClientSettings{
				listFlavorsResp: &mongodbflex.ListFlavorsResponse{
					Flavors: []mongodbflex.InstanceFlavor{
						{
							Id:     new(testFlavorId),
							Cpu:    new(int32(2)),
							Memory: new(int32(4)),
						},
					},
				},
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					FlavorId: new(testFlavorId),
				}),
		},
		{
			description: "update storage class only",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = new("class")
			}),
			isValid: true,
			mockClientSettings: mockClientSettings{
				getInstanceResp: &mongodbflex.InstanceResponse{
					Item: &mongodbflex.Instance{
						Flavor: &mongodbflex.Flavor{
							Id: new(testFlavorId),
						},
					},
				},
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					Storage: &mongodbflex.Storage{
						Class: new("class"),
					},
				}),
		},
		{
			description: "update storage class and size",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = new("class")
				model.StorageSize = new(int64(10))
			}),
			isValid: true,
			mockClientSettings: mockClientSettings{
				getInstanceResp: &mongodbflex.InstanceResponse{
					Item: &mongodbflex.Instance{
						Flavor: &mongodbflex.Flavor{
							Id: new(testFlavorId),
						},
					},
				},
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					Storage: &mongodbflex.Storage{
						Class: new("class"),
						Size:  new(int64(10)),
					},
				}),
		},
		{
			description: "get flavors fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.CPU = new(int32(2))
					model.RAM = new(int32(4))
				},
			),
			mockClientSettings: mockClientSettings{
				listFlavorsFails: true,
			},
			isValid: false,
		},
		{
			description: "flavor id not found",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.CPU = new(int32(5))
					model.RAM = new(int32(9))
				},
			),
			mockClientSettings: mockClientSettings{
				listFlavorsResp: &mongodbflex.ListFlavorsResponse{
					Flavors: []mongodbflex.InstanceFlavor{
						{
							Id:     new(testFlavorId),
							Cpu:    new(int32(2)),
							Memory: new(int32(4)),
						},
						{
							Id:     new("other-flavor"),
							Cpu:    new(int32(1)),
							Memory: new(int32(8)),
						},
					},
				},
			},
			isValid: false,
		},
		{
			description: "get instance fails",
			model:       fixtureRequiredInputModel(),
			mockClientSettings: mockClientSettings{
				getInstanceFails: true,
			},
			expectedRequest: testClient.DefaultAPI.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{}),
			isValid: true,
		},
		{
			description: "get storages fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = new(int32(2))
					model.RAM = new(int32(4))
				},
			),
			mockClientSettings: mockClientSettings{
				listFlavorsFails: true,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, newAPIClientMock(tt.mockClientSettings))
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

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		async         bool
		instanceLabel string
		resp          *mongodbflex.UpdateInstanceResponse
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
			name: "empty response",
			args: args{
				resp: &mongodbflex.UpdateInstanceResponse{},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.args.outputFormat, tt.args.async, tt.args.instanceLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
