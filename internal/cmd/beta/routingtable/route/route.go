package route

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Manage routes of a routing-table",
		Long:  "Manage routes of a routing-table",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
}
