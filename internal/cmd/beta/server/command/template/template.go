package template

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/template/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/template/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Provides functionality for Server Command Template",
		Long:  "Provides functionality for Server Command Template.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
