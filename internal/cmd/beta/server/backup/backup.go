package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup/schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for Server Backup",
		Long:  "Provides functionality for Server Backup.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(disable.NewCmd(p))
	cmd.AddCommand(schedule.NewCmd(p))
}
