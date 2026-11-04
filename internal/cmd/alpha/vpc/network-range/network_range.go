package networkranges

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-range",
		Short: "Provides functionality for network ranges in VPC",
		Long:  "Provides functionality for network ranges in VPC.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
