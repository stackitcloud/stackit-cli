package networkinterface

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/network-interface/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/network-interface/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/network-interface/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-interface",
		Short: "Allows attaching/detaching network interfaces to servers",
		Long:  "Allows attaching/detaching network interfaces to servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(attach.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(detach.NewCmd(params))
}
