package mongodbflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mongodbflex",
		Short: "Provides functionality for MongoDB Flex",
		Long:  "Provides functionality for MongoDB Flex.",
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
