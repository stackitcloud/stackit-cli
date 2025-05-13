package performanceclass

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/performance-class/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/performance-class/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "performance-class",
		Short: "Provides functionality for volume performance classes available inside a project",
		Long:  "Provides functionality for volume performance classes available inside a project.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
