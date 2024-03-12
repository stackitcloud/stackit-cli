package backups

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backups/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backups/list"
	updateschedule "github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/backups/update-schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backups",
		Short: "Provides functionality for PostgreSQL Flex instance backups",
		Long:  "Provides functionality for PostgreSQL Flex instance backups.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(updateschedule.NewCmd())
}
