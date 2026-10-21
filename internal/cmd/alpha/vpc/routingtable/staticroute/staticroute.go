package staticroute

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// The static-route command group is in alpha because the API is only on alpha available.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "static-route",
		Short: "Manages static routes of a VPC routing-table",
		Long:  "Manages static routes of a VPC routing-table",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)

	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
}
