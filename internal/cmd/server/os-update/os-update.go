package osupdate

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os-update",
		Short: "Provides functionality for managed server updates",
		Long:  "Provides functionality for managed server updates.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
	cmd.AddCommand(schedule.NewCmd(params))
}
