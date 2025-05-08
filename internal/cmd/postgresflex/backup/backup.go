package backup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backup/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backup/list"
	updateschedule "github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backup/update-schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Provides functionality for PostgreSQL Flex instance backups",
		Long:  "Provides functionality for PostgreSQL Flex instance backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(updateschedule.NewCmd(params))
}
