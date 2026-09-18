package execution

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/execution/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/execution/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/execution/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution",
		Short: "Provides functionality for Volume Automation Execution",
		Long:  "Provides functionality for Volume Automation Execution.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
