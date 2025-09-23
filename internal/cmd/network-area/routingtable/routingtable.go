package routingtable

import (
	"github.com/spf13/cobra"
	rtCreate "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/create"
	rtDelete "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/delete"
	rtDescribe "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/describe"
	rtList "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/list"
	rtRoute "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/route"
	rtUpdate "github.com/stackitcloud/stackit-cli/internal/cmd/network-area/routingtable/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing-table",
		Short: "Manage routing-tables and its according routes",
		Long: `Manage routing-tables and their associated routes.

This API is currently available only to selected customers.  
To request access, please contact your account manager or submit a support ticket.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(
		rtCreate.NewCmd(params),
		rtUpdate.NewCmd(params),
		rtList.NewCmd(params),
		rtDescribe.NewCmd(params),
		rtDelete.NewCmd(params),
		rtRoute.NewCmd(params),
	)
}
