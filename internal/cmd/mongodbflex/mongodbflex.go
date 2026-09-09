package mongodbflex

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/flavor"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/options"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/storage"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/user"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex/version"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
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

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(user.NewCmd(params))
	cmd.AddCommand(options.NewCmd(params)) //nolint:staticcheck // Command is deprecated but must be kept for backward compatibility
	cmd.AddCommand(backup.NewCmd(params))
	cmd.AddCommand(version.NewCmd(params))
	cmd.AddCommand(flavor.NewCmd(params))
	cmd.AddCommand(storage.NewCmd(params))
}
