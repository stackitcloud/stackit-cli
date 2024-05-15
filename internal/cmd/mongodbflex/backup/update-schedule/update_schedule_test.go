package updateschedule

import (
	"context"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

var projectIdFlag = globalflags.ProjectIdFlag

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{}
var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testSchedule = "0 0/6 * * *"

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		projectIdFlag:  testProjectId,
		scheduleFlag:   testSchedule,
		instanceIdFlag: testInstanceId,
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
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:     utils.Ptr(testInstanceId),
		BackupSchedule: &testSchedule,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixturePayload(mods ...func(payload *mongodbflex.UpdateBackupSchedulePayload)) mongodbflex.UpdateBackupSchedulePayload {
	payload := mongodbflex.UpdateBackupSchedulePayload{
		BackupSchedule:                 utils.Ptr(testSchedule),
		SnapshotRetentionDays:          utils.Ptr(int64(3)),
		DailySnapshotRetentionDays:     utils.Ptr(int64(0)),
		WeeklySnapshotRetentionWeeks:   utils.Ptr(int64(3)),
		MonthlySnapshotRetentionMonths: utils.Ptr(int64(1)),
		PointInTimeWindowHours:         utils.Ptr(int64(30)),
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureUpdateBackupScheduleRequest(mods ...func(request *mongodbflex.ApiUpdateBackupScheduleRequest)) mongodbflex.ApiUpdateBackupScheduleRequest {
	request := testClient.UpdateBackupSchedule(testCtx, testProjectId, testInstanceId)
	request = request.UpdateBackupSchedulePayload(fixturePayload())
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureGetInstanceRequest(mods ...func(request *mongodbflex.ApiGetInstanceRequest)) mongodbflex.ApiGetInstanceRequest {
	request := testClient.GetInstance(testCtx, testProjectId, testInstanceId)
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func fixtureInstance(mods ...func(instance *mongodbflex.Instance)) *mongodbflex.Instance {
	instance := mongodbflex.Instance{
		BackupSchedule: &testSchedule,
		Options: &map[string]string{
			"dailySnapshotRetentionDays":     "0",
			"weeklySnapshotRetentionWeeks":   "3",
			"monthlySnapshotRetentionMonths": "1",
			"pointInTimeWindowHours":         "30",
			"snapshotRetentionDays":          "3",
		},
	}
	for _, mod := range mods {
		mod(&instance)
	}
	return &instance
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
			description: "instance id missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, instanceIdFlag)
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
			description: "backup schedule missing",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, scheduleFlag)
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := NewCmd(nil)
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

			model, err := parseInput(nil, cmd)
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

func TestBuildGetInstanceRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest mongodbflex.ApiGetInstanceRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureGetInstanceRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildGetInstanceRequest(testCtx, tt.model, testClient)

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

func TestBuildUpdateBackupScheduleRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		instance        *mongodbflex.Instance
		expectedRequest mongodbflex.ApiUpdateBackupScheduleRequest
	}{
		{
			description:     "update backup schedule, read retention policy from instance",
			model:           fixtureInputModel(),
			instance:        fixtureInstance(),
			expectedRequest: fixtureUpdateBackupScheduleRequest(),
		},
		{
			description: "update retention policy, read backup schedule from instance",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				InstanceId:                utils.Ptr(testInstanceId),
				DailySnaphotRetentionDays: utils.Ptr(int64(2)),
			},
			instance: fixtureInstance(),
			expectedRequest: fixtureUpdateBackupScheduleRequest().UpdateBackupSchedulePayload(
				fixturePayload(func(payload *mongodbflex.UpdateBackupSchedulePayload) {
					payload.DailySnapshotRetentionDays = utils.Ptr(int64(2))
				}),
			),
		},
		{
			description: "update backup schedule and retention policy",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				InstanceId:                     utils.Ptr(testInstanceId),
				BackupSchedule:                 utils.Ptr("0 0/6 5 2 1"),
				DailySnaphotRetentionDays:      utils.Ptr(int64(2)),
				WeeklySnapshotRetentionWeeks:   utils.Ptr(int64(2)),
				MonthlySnapshotRetentionMonths: utils.Ptr(int64(2)),
				SnapshotRetentionDays:          utils.Ptr(int64(2)),
			},
			instance: fixtureInstance(),
			expectedRequest: fixtureUpdateBackupScheduleRequest().UpdateBackupSchedulePayload(
				fixturePayload(func(payload *mongodbflex.UpdateBackupSchedulePayload) {
					payload.BackupSchedule = utils.Ptr("0 0/6 5 2 1")
					payload.DailySnapshotRetentionDays = utils.Ptr(int64(2))
					payload.WeeklySnapshotRetentionWeeks = utils.Ptr(int64(2))
					payload.MonthlySnapshotRetentionMonths = utils.Ptr(int64(2))
					payload.SnapshotRetentionDays = utils.Ptr(int64(2))
				}),
			),
		},
		{
			description: "no fields set, empty instance (use defaults)",
			model: &inputModel{
				GlobalFlagModel: &globalflags.GlobalFlagModel{
					ProjectId: testProjectId,
				},
				InstanceId: utils.Ptr(testInstanceId),
			},
			instance:        &mongodbflex.Instance{},
			expectedRequest: fixtureUpdateBackupScheduleRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildUpdateBackupScheduleRequest(testCtx, tt.model, tt.instance, testClient)

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
