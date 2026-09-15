package template

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/template/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/template/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Provides functionality for Volume Automation Templates",
		Long:  "Provides functionality for Volume Automation Templates.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
