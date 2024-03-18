package clone

import (
	"context"
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
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testRecoveryTimestamp = "2024-03-08T09:28:00+00:00"

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
	recoveryTimestampString := testRecoveryTimestamp.String()

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
	recoveryTimestampString := testRecoveryTimestamp.String()

	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
		},
		InstanceId:   testInstanceId,
		StorageClass: utils.Ptr("premium-perf4-stackit"),
		StorageSize:  utils.Ptr(int64(10)),
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
	recoveryTimestampString := testRecoveryTimestamp.String()

	payload := postgresflex.CloneInstancePayload{
		Class:     utils.Ptr("premium-perf4-stackit"),
		Size:      utils.Ptr(int64(10)),
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
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest postgresflex.ApiCloneInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureStandardInputModel(),
			expectedRequest: fixtureRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, _ := buildRequest(testCtx, tt.model, testClient)

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
