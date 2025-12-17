package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/restore"
	restorejobs "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/restore-jobs"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/schedule"
	updateschedule "github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup/update-schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for MongoDB Flex instance backups",
		Long:  "Provides functionality for MongoDB Flex instance backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(updateschedule.NewCmd(params))
	cmd.AddCommand(schedule.NewCmd(params))
	cmd.AddCommand(restore.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(restorejobs.NewCmd(params))
}
