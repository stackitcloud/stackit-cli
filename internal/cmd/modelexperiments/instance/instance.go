package instance

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance/get"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance/patch"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for AI Model Experiments instances",
		Long:  "Provides functionality for AI Model Experiments instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}

	addSubcommands(cmd, params)

	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(get.NewCmd(params))
	cmd.AddCommand(patch.NewCmd(params))
}
