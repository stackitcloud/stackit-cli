package network

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Provides functionality for networks",
		Long:  "Provides functionality for networks.",
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
