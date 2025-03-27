package update

import (
	"context"
	"strconv"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &serverbackup.APIClient{}
var testProjectId = uuid.NewString()
var testServerId = uuid.NewString()
var testVolumeId = uuid.NewString()
var testBackupScheduleId = "5"
var testRegion = "eu01"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testBackupScheduleId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		serverIdFlag:              testServerId,
		backupScheduleNameFlag:    "example-backup-schedule-name",
		enabledFlag:               "true",
		rruleFlag:                 defaultRrule,
		backupNameFlag:            "example-backup-name",
		backupRetentionPeriodFlag: "14",
		backupVolumeIdsFlag:       testVolumeId,
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
		BackupScheduleId:      testBackupScheduleId,
		ServerId:              testServerId,
		BackupScheduleName:    utils.Ptr("example-backup-schedule-name"),
		Enabled:               utils.Ptr(defaultEnabled),
		Rrule:                 utils.Ptr(defaultRrule),
		BackupName:            utils.Ptr("example-backup-name"),
		BackupRetentionPeriod: utils.Ptr(int64(14)),
		BackupVolumeIds:       []string{testVolumeId},
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureBackupSchedule(mods ...func(schedule *serverbackup.BackupSchedule)) *serverbackup.BackupSchedule {
	id, _ := strconv.ParseInt(testBackupScheduleId, 10, 64)
	schedule := &serverbackup.BackupSchedule{
		Name:    utils.Ptr("example-backup-schedule-name"),
		Id:      utils.Ptr(id),
		Enabled: utils.Ptr(defaultEnabled),
		Rrule:   utils.Ptr(defaultRrule),
		BackupProperties: &serverbackup.BackupProperties{
			Name:            utils.Ptr("example-backup-name"),
			RetentionPeriod: utils.Ptr(int64(14)),
			VolumeIds:       utils.Ptr([]string{testVolumeId}),
		},
	}
	for _, mod := range mods {
		mod(schedule)
	}
	return schedule
}

func fixturePayload(mods ...func(payload *serverbackup.UpdateBackupSchedulePayload)) serverbackup.UpdateBackupSchedulePayload {
	payload := serverbackup.UpdateBackupSchedulePayload{
		Name:    utils.Ptr("example-backup-schedule-name"),
		Enabled: utils.Ptr(defaultEnabled),
		Rrule:   utils.Ptr("DTSTART;TZID=Europe/Sofia:20200803T023000 RRULE:FREQ=DAILY;INTERVAL=1"),
		BackupProperties: &serverbackup.BackupProperties{
			Name:            utils.Ptr("example-backup-name"),
			RetentionPeriod: utils.Ptr(int64(14)),
			VolumeIds:       utils.Ptr([]string{testVolumeId}),
		},
	}
	for _, mod := range mods {
		mod(&payload)
	}
	return payload
}

func fixtureRequest(mods ...func(request *serverbackup.ApiUpdateBackupScheduleRequest)) serverbackup.ApiUpdateBackupScheduleRequest {
	request := testClient.UpdateBackupSchedule(testCtx, testProjectId, testServerId, testRegion, testBackupScheduleId)
	request = request.UpdateBackupSchedulePayload(fixturePayload())
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
		isValid       bool
		expectedModel *inputModel
	}{

		{
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "backup schedule id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
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

			err = cmd.ValidateFlagGroups()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flag groups: %v", err)
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
		description     string
		model           *inputModel
		expectedRequest serverbackup.ApiUpdateBackupScheduleRequest
		isValid         bool
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
			isValid:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(testCtx, tt.model, testClient, *fixtureBackupSchedule())
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
		outputFormat string
		resp         serverbackup.BackupSchedule
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(p)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
