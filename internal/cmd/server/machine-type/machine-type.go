package machinetype

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/machine-type/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/machine-type/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machine-type",
		Short: "Provides functionality for server machine types available inside a project",
		Long:  "Provides functionality for server machine types available inside a project.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
}
