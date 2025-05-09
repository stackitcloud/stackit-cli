package logme

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/plans"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logme",
		Short: "Provides functionality for LogMe",
		Long:  "Provides functionality for LogMe.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
}
