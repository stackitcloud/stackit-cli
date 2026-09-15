package automation

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/template"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "automation",
		Short: "Provides functionality for Volume Automation",
		Long:  "Provides functionality for Volume Automation.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(template.NewCmd(params))
}
