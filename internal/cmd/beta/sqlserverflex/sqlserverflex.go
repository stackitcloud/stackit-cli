package sqlserverflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/database"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/user"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sqlserverflex",
		Short: "Provides functionality for SQLServer Flex",
		Long:  "Provides functionality for SQLServer Flex.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(database.NewCmd(params))
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params))
	cmd.AddCommand(user.NewCmd(params))
}
