package networkarea

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/list"
	networkrange "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/region"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-area",
		Short: "Provides functionality for STACKIT Network Area (SNA)",
		Long:  "Provides functionality for STACKIT Network Area (SNA).",
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
	cmd.AddCommand(networkrange.NewCmd(params))
	cmd.AddCommand(region.NewCmd(params))
	cmd.AddCommand(route.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
