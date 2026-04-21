package utils

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

type ServerBackupClient interface {
	ListBackupSchedulesExecute(ctx context.Context, projectId, serverId, region string) (*serverbackup.GetBackupSchedulesResponse, error)
	ListBackupsExecute(ctx context.Context, projectId, serverId, region string) (*serverbackup.GetBackupsListResponse, error)
}

func CanDisableBackupService(ctx context.Context, client ServerBackupClient, projectId, serverId, region string) (bool, error) {
	schedules, err := client.ListBackupSchedulesExecute(ctx, projectId, serverId, region)
	if err != nil {
		return false, fmt.Errorf("list backup schedules: %w", err)
	}
	if len(*schedules.Items) > 0 {
		return false, nil
	}

	backups, err := client.ListBackupsExecute(ctx, projectId, serverId, region)
	if err != nil {
		return false, fmt.Errorf("list backups: %w", err)
	}
	if len(*backups.Items) > 0 {
		return false, nil
	}

	// no backups and no backup schedules found for this server => can disable backup service
	return true, nil
}
