package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for MariaDB credentials",
		Long:  "Provides functionality for MariaDB credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
}
