package flavor

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/flavor/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex/flavor/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flavor",
		Short: "Provides functionality for SQLServer Flex flavors",
		Long:  "Provides functionality for SQLServer Flex flavors.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
