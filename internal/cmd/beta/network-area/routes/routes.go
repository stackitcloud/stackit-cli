package routes

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/routes/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/routes/delete"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Provides functionality for static routes in STACKIT Network Areas",
		Long:  "Provides functionality for static routes in STACKIT Network Areas.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
}
