package networkinterface

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-interface/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-interface/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-interface/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-interface/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-interface/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-interface",
		Short: "Provides functionality for network interfaces",
		Long:  "Provides functionality for network interfaces.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
