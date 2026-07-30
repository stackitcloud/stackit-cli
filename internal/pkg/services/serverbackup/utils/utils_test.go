package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/uuid"
	serverbackup "github.com/stackitcloud/stackit-sdk-go/services/serverbackup/v2api"
)

const testRegion = "eu01"

var (
	testProjectId = uuid.NewString()
	testServerId  = uuid.NewString()
)

type mockSettings struct {
	listBackupSchedulesFails bool
	listBackupSchedulesResp  *serverbackup.GetBackupSchedulesResponse
	listBackupsFails         bool
	listBackupsResp          *serverbackup.GetBackupsListResponse
}

func newServerbackupClientMock(s mockSettings) serverbackup.DefaultAPI {
	return &serverbackup.DefaultAPIServiceMock{
		ListBackupSchedulesExecuteMock: utils.Ptr(func(_ serverbackup.ApiListBackupSchedulesRequest) (*serverbackup.GetBackupSchedulesResponse, error) {
			if s.listBackupSchedulesFails {
				return nil, fmt.Errorf("could not list backup schedules")
			}
			return s.listBackupSchedulesResp, nil
		}),
		ListBackupsExecuteMock: utils.Ptr(func(_ serverbackup.ApiListBackupsRequest) (*serverbackup.GetBackupsListResponse, error) {
			if s.listBackupsFails {
				return nil, fmt.Errorf("could not list backups")
			}
			return s.listBackupsResp, nil
		}),
	}
}

func TestCanDisableBackupService(t *testing.T) {
	tests := []struct {
		description    string
		mockSettings   mockSettings
		isValid        bool // isValid ==> err == nil
		expectedOutput bool // expectedCanDisable
	}{
		{
			description: "base-ok-can-disable-backups-service-no-backups-no-backup-schedules",
			mockSettings: mockSettings{
				listBackupsFails:         false,
				listBackupSchedulesFails: false,
				listBackupsResp:          &serverbackup.GetBackupsListResponse{Items: []serverbackup.Backup{}},
				listBackupSchedulesResp:  &serverbackup.GetBackupSchedulesResponse{Items: []serverbackup.BackupSchedule{}},
			},
			isValid:        true,
			expectedOutput: true,
		},
		{
			description: "not-ok-api-error-list-backups",
			mockSettings: mockSettings{
				listBackupsFails:         true,
				listBackupSchedulesFails: false,
				listBackupsResp:          &serverbackup.GetBackupsListResponse{Items: []serverbackup.Backup{}},
				listBackupSchedulesResp:  &serverbackup.GetBackupSchedulesResponse{Items: []serverbackup.BackupSchedule{}},
			},
			isValid:        false,
			expectedOutput: false,
		},
		{
			description: "not-ok-api-error-list-backup-schedules",
			mockSettings: mockSettings{
				listBackupsFails:         true,
				listBackupSchedulesFails: false,
				listBackupsResp:          &serverbackup.GetBackupsListResponse{Items: []serverbackup.Backup{}},
				listBackupSchedulesResp:  &serverbackup.GetBackupSchedulesResponse{Items: []serverbackup.BackupSchedule{}},
			},
			isValid:        false,
			expectedOutput: false,
		},
		{
			description: "not-ok-has-backups-cannot-disable",
			mockSettings: mockSettings{
				listBackupsFails:         false,
				listBackupSchedulesFails: false,
				listBackupsResp: &serverbackup.GetBackupsListResponse{
					Items: []serverbackup.Backup{
						{
							CreatedAt:      "test timestamp",
							ExpireAt:       "test timestamp",
							Id:             "5",
							LastRestoredAt: utils.Ptr("test timestamp"),
							Name:           "test name",
							Size:           utils.Ptr(int32(5)),
							Status:         serverbackup.BACKUPSTATUS_IN_PROGRESS,
							VolumeBackups:  nil,
						},
					},
				},
				listBackupSchedulesResp: &serverbackup.GetBackupSchedulesResponse{Items: []serverbackup.BackupSchedule{}},
			},
			isValid:        true,
			expectedOutput: false,
		},
		{
			description: "not-ok-has-backups-schedules-cannot-disable",
			mockSettings: mockSettings{
				listBackupsFails:         false,
				listBackupSchedulesFails: false,
				listBackupsResp:          &serverbackup.GetBackupsListResponse{Items: []serverbackup.Backup{}},
				listBackupSchedulesResp: &serverbackup.GetBackupSchedulesResponse{
					Items: []serverbackup.BackupSchedule{
						{
							BackupProperties: nil,
							Enabled:          false,
							Id:               int32(5),
							Name:             "some name",
							Rrule:            "some rrule",
						},
					},
				},
			},
			isValid:        true,
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newServerbackupClientMock(tt.mockSettings)

			output, err := CanDisableBackupService(context.Background(), client, testProjectId, testServerId, testRegion)

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
