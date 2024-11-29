package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/create"
	del "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/restore"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/schedule"
	volumebackup "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/volume-backup"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for server backups",
		Long:  "Provides functionality for server backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(disable.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(schedule.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(restore.NewCmd(p))
	cmd.AddCommand(del.NewCmd(p))
	cmd.AddCommand(volumebackup.NewCmd(p))
}
