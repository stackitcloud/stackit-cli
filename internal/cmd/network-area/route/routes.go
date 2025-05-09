package route

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network-area/route/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route",
		Short: "Provides functionality for static routes in STACKIT Network Areas",
		Long:  "Provides functionality for static routes in STACKIT Network Areas.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
