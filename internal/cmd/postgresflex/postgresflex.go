package postgresflex

import (
	"stackit/internal/cmd/postgresflex/instance"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "postgresflex",
		Short: "Provides functionality for PostgreSQL Flex",
		Long:  "Provides functionality for PostgreSQL Flex",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
}
