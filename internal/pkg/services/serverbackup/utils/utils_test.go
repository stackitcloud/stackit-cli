package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

var (
	testProjectId = uuid.NewString()
	testServerId  = uuid.NewString()
)

type serverbackupClientMocked struct {
	listBackupSchedulesFails bool
	listBackupSchedulesResp  *serverbackup.GetBackupSchedulesResponse
	listBackupsFails         bool
	listBackupsResp          *serverbackup.GetBackupsListResponse
}

func (m *serverbackupClientMocked) ListBackupSchedulesExecute(_ context.Context, _, _ string) (*serverbackup.GetBackupSchedulesResponse, error) {
	if m.listBackupSchedulesFails {
		return nil, fmt.Errorf("could not list backup schedules")
	}
	return m.listBackupSchedulesResp, nil
}

func (m *serverbackupClientMocked) ListBackupsExecute(_ context.Context, _, _ string) (*serverbackup.GetBackupsListResponse, error) {
	if m.listBackupsFails {
		return nil, fmt.Errorf("could not list backups")
	}
	return m.listBackupsResp, nil
}

func TestCanDisableBackupService(t *testing.T) {
	tests := []struct {
		description              string
		listBackupsFails         bool
		listBackupSchedulesFails bool
		listBackups              *serverbackup.GetBackupsListResponse
		listBackupSchedules      *serverbackup.GetBackupSchedulesResponse
		isValid                  bool // isValid ==> err == nil
		expectedOutput           bool // expectedCanDisable
	}{
		{
			description:              "base-ok-can-disable-backups-service-no-backups-no-backup-schedules",
			listBackupsFails:         false,
			listBackupSchedulesFails: false,
			listBackups:              &serverbackup.GetBackupsListResponse{Items: &[]serverbackup.Backup{}},
			listBackupSchedules:      &serverbackup.GetBackupSchedulesResponse{Items: &[]serverbackup.BackupSchedule{}},
			isValid:                  true,
			expectedOutput:           true,
		},
		{
			description:              "not-ok-api-error-list-backups",
			listBackupsFails:         true,
			listBackupSchedulesFails: false,
			listBackups:              &serverbackup.GetBackupsListResponse{Items: &[]serverbackup.Backup{}},
			listBackupSchedules:      &serverbackup.GetBackupSchedulesResponse{Items: &[]serverbackup.BackupSchedule{}},
			isValid:                  false,
			expectedOutput:           false,
		},
		{
			description:              "not-ok-api-error-list-backup-schedules",
			listBackupsFails:         true,
			listBackupSchedulesFails: false,
			listBackups:              &serverbackup.GetBackupsListResponse{Items: &[]serverbackup.Backup{}},
			listBackupSchedules:      &serverbackup.GetBackupSchedulesResponse{Items: &[]serverbackup.BackupSchedule{}},
			isValid:                  false,
			expectedOutput:           false,
		},
		{
			description:              "not-ok-has-backups-cannot-disable",
			listBackupsFails:         false,
			listBackupSchedulesFails: false,
			listBackups: &serverbackup.GetBackupsListResponse{
				Items: &[]serverbackup.Backup{
					{
						CreatedAt:      utils.Ptr("test timestamp"),
						ExpireAt:       utils.Ptr("test timestamp"),
						Id:             utils.Ptr("5"),
						LastRestoredAt: utils.Ptr("test timestamp"),
						Name:           utils.Ptr("test name"),
						Size:           utils.Ptr(int64(5)),
						Status:         utils.Ptr("test status"),
						VolumeBackups:  nil,
					},
				},
			},
			listBackupSchedules: &serverbackup.GetBackupSchedulesResponse{Items: &[]serverbackup.BackupSchedule{}},
			isValid:             true,
			expectedOutput:      false,
		},
		{
			description:              "not-ok-has-backups-schedules-cannot-disable",
			listBackupsFails:         false,
			listBackupSchedulesFails: false,
			listBackups:              &serverbackup.GetBackupsListResponse{Items: &[]serverbackup.Backup{}},
			listBackupSchedules: &serverbackup.GetBackupSchedulesResponse{
				Items: &[]serverbackup.BackupSchedule{
					{
						BackupProperties: nil,
						Enabled:          utils.Ptr(false),
						Id:               utils.Ptr(int64(5)),
						Name:             utils.Ptr("some name"),
						Rrule:            utils.Ptr("some rrule"),
					},
				},
			},
			isValid:        true,
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &serverbackupClientMocked{
				listBackupsFails:         tt.listBackupsFails,
				listBackupSchedulesFails: tt.listBackupSchedulesFails,
				listBackupsResp:          tt.listBackups,
				listBackupSchedulesResp:  tt.listBackupSchedules,
			}

			output, err := CanDisableBackupService(context.Background(), client, testProjectId, testServerId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %t, got %t", tt.expectedOutput, output)
			}
		})
	}
}
