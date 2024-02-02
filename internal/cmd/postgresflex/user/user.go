package user

import (
	"stackit/internal/cmd/postgresflex/user/create"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for PostgreSQL Flex users",
		Long:  "Provides functionality for PostgreSQL Flex users",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
}
