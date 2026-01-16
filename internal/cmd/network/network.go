package network

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Provides functionality for networks",
		Long:  "Provides functionality for networks.",
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
	cmd.AddCommand(update.NewCmd(params))
}
