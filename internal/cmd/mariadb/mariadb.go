package mariadb

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mariadb",
		Short: "Provides functionality for MariaDB",
		Long:  "Provides functionality for MariaDB.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(plans.NewCmd())
	cmd.AddCommand(credentials.NewCmd())
}
