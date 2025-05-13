package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/create"
	del "github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/restore"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/schedule"
	volumebackup "github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/volume-backup"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for server backups",
		Long:  "Provides functionality for server backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(schedule.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
	cmd.AddCommand(del.NewCmd(params))
	cmd.AddCommand(volumebackup.NewCmd(params))
}
