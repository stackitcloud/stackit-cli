package performanceclass

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/performance-class/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/performance-class/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "performance-class",
		Short: "Provides functionality for volume performance classes available inside a project",
		Long:  "Provides functionality for volume performance classes available inside a project.",
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
