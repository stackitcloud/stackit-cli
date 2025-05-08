package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	completerotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/complete-rotation"
	startrotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/start-rotation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for SKE credentials",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(startrotation.NewCmd(params))
	cmd.AddCommand(completerotation.NewCmd(params))
}
