package network

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Provides functionality for Network",
		Long:  "Provides functionality for Network.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
