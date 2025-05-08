package volumebackup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	del "github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/volume-backup/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup/volume-backup/restore"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume-backup",
		Short: "Provides functionality for Server Backup Volume Backups",
		Long:  "Provides functionality for Server Backup Volume Backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(del.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
}
