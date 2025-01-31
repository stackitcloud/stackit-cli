package machinetype

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/machine-type/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/machine-type/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machine-type",
		Short: "Provides functionality for server machine types available inside a project",
		Long:  "Provides functionality for server machine types available inside a project.",
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
