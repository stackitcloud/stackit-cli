package clone

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
	listStoragesFails bool
	listStoragesResp  *postgresflex.ListStoragesResponse
	getInstanceFails  bool
	getInstanceResp   *postgresflex.InstanceResponse
}

func (c *postgresFlexClientMocked) CloneInstance(ctx context.Context, projectId, instanceId string) postgresflex.ApiCloneInstanceRequest {
	return testClient.CloneInstance(ctx, projectId, instanceId)
}

func (c *postgresFlexClientMocked) GetInstanceExecute(_ context.Context, _, _ string) (*postgresflex.InstanceResponse, error) {
	if c.getInstanceFails {
		return nil, fmt.Errorf("get instance failed")
	}
	return c.getInstanceResp, nil
}

func (c *postgresFlexClientMocked) ListStoragesExecute(_ context.Context, _, _ string) (*postgresflex.ListStoragesResponse, error) {
	if c.listStoragesFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return c.listStoragesResp, nil
}

var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testRecoveryTimestamp = "2024-03-08T09:28:00+00:00"
var testFlavorId = uuid.NewString()
var testStorageClass = "premium-perf4-stackit"
var testStorageSize = int64(10)

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
		projectIdFlag:         testProjectId,
		recoveryTimestampFlag: testRecoveryTimestamp,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureStandardFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:         testProjectId,
		recoveryTimestampFlag: testRecoveryTimestamp,
		storageClassFlag:      "class",
		storageSizeFlag:       "10",
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
	recoveryTimestampString := testRecoveryTimestamp.Format(time.RFC3339)

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
		},
		InstanceId:   testInstanceId,
		RecoveryDate: utils.Ptr(recoveryTimestampString),
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
	recoveryTimestampString := testRecoveryTimestamp.Format(time.RFC3339)

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
		},
		InstanceId:   testInstanceId,
		StorageClass: utils.Ptr(testStorageClass),
		StorageSize:  utils.Ptr(testStorageSize),
		RecoveryDate: utils.Ptr(recoveryTimestampString),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *postgresflex.ApiCloneInstanceRequest)) postgresflex.ApiCloneInstanceRequest {
	request := testClient.CloneInstance(testCtx, testProjectId, testInstanceId)
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
	recoveryTimestampString := testRecoveryTimestamp.Format(time.RFC3339)

	payload := postgresflex.CloneInstancePayload{
		Timestamp: utils.Ptr(recoveryTimestampString),
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

			model, err := parseInput(cmd, tt.argValues)
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
	testRecoveryTimestamp, err := time.Parse(recoveryDateFormat, testRecoveryTimestamp)
	if err != nil {
		return
	}
	recoveryTimestampString := testRecoveryTimestamp.Format(time.RFC3339)

	tests := []struct {
		description       string
		model             *inputModel
		expectedRequest   postgresflex.ApiCloneInstanceRequest
		getInstanceFails  bool
		getInstanceResp   *postgresflex.InstanceResponse
		listStoragesFails bool
		listStoragesResp  *postgresflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "base",
			model:           fixtureRequiredInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "specify storage class only",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr("class")
			}),
			isValid: true,
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Flavor: &postgresflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
					Storage: &postgresflex.Storage{
						Class: utils.Ptr(testStorageClass),
						Size:  utils.Ptr(testStorageSize),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			expectedRequest: testClient.CloneInstance(testCtx, testProjectId, testInstanceId).
				CloneInstancePayload(postgresflex.CloneInstancePayload{
					Class:     utils.Ptr("class"),
					Timestamp: utils.Ptr(recoveryTimestampString),
				}),
		},
		{
			description: "specify storage class and size",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr("class")
				model.StorageSize = utils.Ptr(int64(10))
			}),
			isValid: true,
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Flavor: &postgresflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
					Storage: &postgresflex.Storage{
						Class: utils.Ptr(testStorageClass),
						Size:  utils.Ptr(testStorageSize),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			expectedRequest: testClient.CloneInstance(testCtx, testProjectId, testInstanceId).
				CloneInstancePayload(postgresflex.CloneInstancePayload{
					Class:     utils.Ptr("class"),
					Size:      utils.Ptr(int64(10)),
					Timestamp: utils.Ptr(recoveryTimestampString),
				}),
		},
		{
			description: "get instance fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("class")
					model.RecoveryDate = utils.Ptr(recoveryTimestampString)
				},
			),
			getInstanceFails: true,
			isValid:          false,
		},
		{
			description: "invalid storage class",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("non-existing-class")
				},
			),
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Flavor: &postgresflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
					Storage: &postgresflex.Storage{
						Class: utils.Ptr(testStorageClass),
						Size:  utils.Ptr(testStorageSize),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &postgresflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
		{
			description: "invalid storage size",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageSize = utils.Ptr(int64(9))
				},
			),
			getInstanceResp: &postgresflex.InstanceResponse{
				Item: &postgresflex.Instance{
					Flavor: &postgresflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
					Storage: &postgresflex.Storage{
						Class: utils.Ptr(testStorageClass),
						Size:  utils.Ptr(testStorageSize),
					},
				},
			},
			listStoragesResp: &postgresflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
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
				getInstanceFails:  tt.getInstanceFails,
				getInstanceResp:   tt.getInstanceResp,
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
