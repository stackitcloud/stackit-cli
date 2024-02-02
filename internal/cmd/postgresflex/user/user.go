package user

import (
	"stackit/internal/cmd/postgresflex/user/create"
	"stackit/internal/cmd/postgresflex/user/delete"
	"stackit/internal/cmd/postgresflex/user/describe"
	"stackit/internal/cmd/postgresflex/user/list"
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
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(delete.NewCmd())
}
