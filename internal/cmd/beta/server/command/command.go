package command

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command/template"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command",
		Short: "Provides functionality for Server Command",
		Long:  "Provides functionality for Server Command.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(template.NewCmd(p))
}
