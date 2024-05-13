package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/list"
	restorejobs "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/restore-jobs"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/schedule"
	updateschedule "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/update-schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for MongoDB Flex instance backups",
		Long:  "Provides functionality for MongoDB Flex instance backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(restorejobs.NewCmd(p))
	cmd.AddCommand(updateschedule.NewCmd(p))
	cmd.AddCommand(schedule.NewCmd(p))
}
