package networkranges

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/network-ranges/create"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network-ranges",
		Aliases: []string{"ranges"},
		Short:   "Provides functionality for Network Area",
		Long:    "Provides functionality for Network Area.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
}
