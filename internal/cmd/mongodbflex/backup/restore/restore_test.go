package restore

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

type testCtxKey struct{}

const (
	testRegion    = "eu02"
	testBackupId  = "backupID"
	testTimestamp = "2021-01-01T00:00:00Z"
)

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{}

var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testBackupInstanceId = uuid.NewString()

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		backupIdFlag:              testBackupId,
		backupInstanceIdFlag:      testBackupInstanceId,
		instanceIdFlag:            testInstanceId,
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
		InstanceId:       testInstanceId,
		BackupId:         testBackupId,
		BackupInstanceId: testBackupInstanceId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRestoreRequest(mods ...func(request mongodbflex.ApiRestoreInstanceRequest)) mongodbflex.ApiRestoreInstanceRequest {
	request := testClient.RestoreInstance(testCtx, testProjectId, testInstanceId, testRegion)
	request = request.RestoreInstancePayload(mongodbflex.RestoreInstancePayload{
		BackupId:   utils.Ptr(testBackupId),
		InstanceId: utils.Ptr(testBackupInstanceId),
	})
	for _, mod := range mods {
		mod(request)
	}
	return request
}

func fixtureCloneRequest(mods ...func(request mongodbflex.ApiCloneInstanceRequest)) mongodbflex.ApiCloneInstanceRequest {
	request := testClient.CloneInstance(testCtx, testProjectId, testInstanceId, testRegion)
	request = request.CloneInstancePayload(mongodbflex.CloneInstancePayload{
		Timestamp:  utils.Ptr(testTimestamp),
		InstanceId: utils.Ptr(testBackupInstanceId),
	})
	for _, mod := range mods {
		mod(request)
	}
	return request
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
			description: "no values",
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
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
			description: "instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[instanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "backup instance id invalid 1",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[backupInstanceIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "backup instance id invalid 2",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[backupInstanceIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "timestamp and backup id both provided",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[timestampFlag] = testTimestamp
			}),
			isValid: false,
		},
		{
			description: "timestamp and backup id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, backupIdFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
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

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			err = cmd.ValidateFlagGroups()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
			}

			model, err := parseInput(p, cmd)
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

func TestBuildRestoreRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest mongodbflex.ApiRestoreInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRestoreRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRestoreRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx))
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildCloneRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest mongodbflex.ApiCloneInstanceRequest
	}{
		{
			description: "base",
			model: fixtureInputModel(func(model *inputModel) {
				model.BackupId = ""
				model.Timestamp = testTimestamp
			}),
			expectedRequest: fixtureCloneRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildCloneRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx))
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestGetIsRestoreOperation(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		expected    bool
	}{
		{
			description: "true",
			model:       fixtureInputModel(),
			expected:    true,
		},
		{
			description: "false",
			model: fixtureInputModel(func(model *inputModel) {
				model.BackupId = ""
				model.Timestamp = testTimestamp
			}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := getIsRestoreOperation(tt.model.BackupId, tt.model.Timestamp)
			if result != tt.expected {
				t.Fatalf("Data does not match: %t", result)
			}
		})
	}
}
