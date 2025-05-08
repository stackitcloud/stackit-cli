package schedule

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/create"
	del "github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Provides functionality for Server os-update Schedule",
		Long:  "Provides functionality for Server os-update Schedule.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(del.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
