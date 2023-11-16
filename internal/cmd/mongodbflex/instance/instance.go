package instance

import (
	"stackit/internal/cmd/mongodbflex/instance/create"
	"stackit/internal/cmd/mongodbflex/instance/delete"
	"stackit/internal/cmd/mongodbflex/instance/describe"
	"stackit/internal/cmd/mongodbflex/instance/list"
	"stackit/internal/cmd/mongodbflex/instance/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for MongoDB Flex instances",
		Long:  "Provides functionality for MongoDB Flex instances",
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
	cmd.AddCommand(update.NewCmd())
}
