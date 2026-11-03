package platform

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/platform/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/platform/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/platform/quota"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit beta scf platform' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Provides functionality for STACKIT Cloud Foundry platforms",
		Long:  "Provides functionality for STACKIT Cloud Foundry platforms.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(quota.NewCmd(params))
}
