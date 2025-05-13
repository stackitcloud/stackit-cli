package networkranges

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network-range",
		Aliases: []string{"range"},
		Short:   "Provides functionality for network ranges in STACKIT Network Areas",
		Long:    "Provides functionality for network ranges in STACKIT Network Areas.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
