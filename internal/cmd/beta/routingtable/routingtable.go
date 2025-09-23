package routingtable

import (
	"github.com/spf13/cobra"
	rtDescribe "github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/describe"
	rtList "github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/list"
	route "github.com/stackitcloud/stackit-cli/internal/cmd/beta/routingtable/route"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing-table",
		Short: "Manage routing-tables and its according routes",
		Long: `Manage routing tables and their associated routes.

This functionality is currently in BETA. At this stage, only listing and describing
routing-tables, as well as full CRUD operations for routes, are supported. 
This feature is primarily intended for debugging routes created through Terraform.

Once the feature reaches General Availability, we plan to introduce support
for creating routing tables and attaching them to networks directly via the
CLI. Until then, we recommend users continue managing routing tables and 
attachments through the Terraform provider.`,
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
