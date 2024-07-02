package volumebackup

import (
	del "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/volume-backup/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/volume-backup/restore"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume-backup",
		Short: "Provides functionality for Server Backup Volume Backups",
		Long:  "Provides functionality for Server Backup Volume Backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(del.NewCmd(p))
	cmd.AddCommand(restore.NewCmd(p))
}
