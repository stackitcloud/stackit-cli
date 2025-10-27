package routingtable

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	rtDescribe "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/describe"
	rtList "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/route"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing-table",
		Short: "Manage routing-tables and its according routes",
		Long: `Manage routing tables and their associated routes.

This API is currently in a private alpha stage. To request access,
please contact your Account Manager or open a support ticket.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		rtList.NewCmd(params),
		rtDescribe.NewCmd(params),
		route.NewCmd(params),
	)
}
