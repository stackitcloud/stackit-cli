package sqlserverflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/database"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/flavor"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/user"
	"github.com/stackitcloud/stackit-cli/internal/cmd/sqlserverflex/version"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
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

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(database.NewCmd(params))
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params)) //nolint:staticcheck // Command is deprecated but must be kept for backward compatibility
	cmd.AddCommand(user.NewCmd(params))
	cmd.AddCommand(version.NewCmd(params))
	cmd.AddCommand(flavor.NewCmd(params))
}
