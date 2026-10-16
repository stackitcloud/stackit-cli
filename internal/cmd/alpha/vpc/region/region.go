package region

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/region/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/region/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/region/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/region/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc/region/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "Manages regional configurations of a VPC",
		Long:  "Manages regional configurations of a VPC.",
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
}
