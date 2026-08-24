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

// Deprecated: Will be removed after 2027-08-31.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:        "redis",
		Deprecated: "\nRedis commands `stackit redis` have been deprecated and will be removed after 31.08.2027. Please use valkey (Key Value Store) commands `stackit valkey` instead.",
		Short:      "Provides functionality for Redis",
		Long:       "Provides functionality for Redis.",
		Args:       args.NoArgs,
		Run:        utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

// Deprecated: Will be removed after 2027-08-31.
func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
}
