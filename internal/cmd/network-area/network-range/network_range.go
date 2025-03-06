package networkranges

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/network-range/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network-range",
		Aliases: []string{"range"},
		Short:   "Provides functionality for network ranges in STACKIT Network Areas",
		Long:    "Provides functionality for network ranges in STACKIT Network Areas.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
