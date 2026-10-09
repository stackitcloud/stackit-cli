package vpc

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/list"
	networkRange "github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/network-range"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpc",
		Short: "Manages vpcs",
		Long:  "Manages the lifecycle of vpcs.",
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
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(networkRange.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
