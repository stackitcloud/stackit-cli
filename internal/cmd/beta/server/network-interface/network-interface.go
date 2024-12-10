package networkinterface

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/network-interface/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/network-interface/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/network-interface/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-interface",
		Short: "Allows attaching/detaching network interfaces to servers",
		Long:  "Allows attaching/detaching network interfaces to servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(attach.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(detach.NewCmd(p))
}
