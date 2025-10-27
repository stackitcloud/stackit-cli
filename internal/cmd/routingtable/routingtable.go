package routingtable

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	rtCreate "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/create"
	rtDelete "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/delete"
	rtDescribe "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/describe"
	rtList "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/route"
	rtUpdate "github.com/stackitcloud/stackit-cli/internal/cmd/routingtable/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
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

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		rtCreate.NewCmd(params),
		rtUpdate.NewCmd(params),
		rtList.NewCmd(params),
		rtDescribe.NewCmd(params),
		rtDelete.NewCmd(params),
		route.NewCmd(params),
	)
}
