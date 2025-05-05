package git

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Provides functionality for STACKIT Git",
		Long:  "Provides functionality for STACKIT Git.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		list.NewCmd(p),
		describe.NewCmd(p),
		create.NewCmd(p),
		delete.NewCmd(p),
	)
}
