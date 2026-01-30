package logs

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/instance"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Provides functionality for Logs",
		Long:  "Provides functionality for Logs.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(access_token.NewCmd(params))
}
