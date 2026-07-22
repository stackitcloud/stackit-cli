package utils

import (
	"context"
	"fmt"

	serverbackup "github.com/stackitcloud/stackit-sdk-go/services/serverbackup/v2api"
)

func CanDisableBackupService(ctx context.Context, client serverbackup.DefaultAPI, projectId, serverId, region string) (bool, error) {
	schedules, err := client.ListBackupSchedules(ctx, projectId, serverId, region).Execute()
	if err != nil {
		return false, fmt.Errorf("list backup schedules: %w", err)
	}
	if len(schedules.Items) > 0 {
		return false, nil
	}

	backups, err := client.ListBackups(ctx, projectId, serverId, region).Execute()
	if err != nil {
		return false, fmt.Errorf("list backups: %w", err)
	}
	if len(backups.Items) > 0 {
		return false, nil
	}

	// no backups and no backup schedules found for this server => can disable backup service
	return true, nil
}
