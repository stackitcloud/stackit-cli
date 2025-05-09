package mongodbflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mongodbflex",
		Short: "Provides functionality for MongoDB Flex",
		Long:  "Provides functionality for MongoDB Flex.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(user.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params))
	cmd.AddCommand(backup.NewCmd(params))
}
