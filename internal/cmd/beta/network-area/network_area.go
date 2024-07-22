package networkarea

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/network-area/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-area",
		Short: "Provides functionality for Network Area",
		Long:  "Provides functionality for Network Area.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
}
