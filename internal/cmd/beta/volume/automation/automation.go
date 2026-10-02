package automation

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/execution"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/template"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/volume/automation/update"
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
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(execution.NewCmd(params))
	cmd.AddCommand(generatepayload.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(template.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
