package redis

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redis",
		Short: "Provides functionality for Redis",
		Long:  "Provides functionality for Redis.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
}
