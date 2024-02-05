package postgresflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresflex",
		Aliases: []string{"postgresqlflex"},
		Short:   "Provides functionality for PostgreSQL Flex",
		Long:    "Provides functionality for PostgreSQL Flex.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(user.NewCmd())
}
