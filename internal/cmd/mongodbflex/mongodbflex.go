package mongodbflex

import (
	"stackit/internal/cmd/mongodbflex/instance"
	"stackit/internal/cmd/mongodbflex/options"
	"stackit/internal/cmd/mongodbflex/user"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mongodbflex",
		Short: "Provides functionality for MongoDB Flex",
		Long:  "Provides functionality for MongoDB Flex",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(user.NewCmd())
	cmd.AddCommand(options.NewCmd())
}
