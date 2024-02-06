package redis

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redis",
		Short: "Provides functionality for Redis",
		Long:  "Provides functionality for Redis.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(plans.NewCmd())
	cmd.AddCommand(credentials.NewCmd())
}
