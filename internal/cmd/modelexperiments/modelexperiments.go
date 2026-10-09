package modelexperiments

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai-model-experiments",
		Short: "Provides functionality for AI Model Experiments",
		Long:  "Provides functionality for AI Model Experiments.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}

	addSubcommands(cmd, params)

	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(token.NewCmd(params))
}
