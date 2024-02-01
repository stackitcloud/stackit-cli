package instance

import (
	"stackit/internal/cmd/postgresflex/instance/create"
	"stackit/internal/cmd/postgresflex/instance/describe"
	"stackit/internal/cmd/postgresflex/instance/list"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for PostgreSQL Flex instances",
		Long:  "Provides functionality for PostgreSQL Flex instances",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(describe.NewCmd())
}
