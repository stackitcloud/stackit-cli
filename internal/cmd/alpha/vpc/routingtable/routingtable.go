package routingtable

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/staticroute"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/routingtable/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing-table",
		Short: "Manages VPC routing-tables",
		Long:  "Manages routing-tables for a VPC in a region.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}

	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(staticroute.NewCmd(params))
}
