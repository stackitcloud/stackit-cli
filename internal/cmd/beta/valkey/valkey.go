package valkey

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valkey",
		Short: "Provides functionality for Valkey",
		Long:  "Provides functionality for Valkey.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(credentials.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(instance.NewCmd(params))
}
