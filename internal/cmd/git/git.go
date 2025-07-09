package git

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/flavor"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Provides functionality for STACKIT Git",
		Long:  "Provides functionality for STACKIT Git.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		instance.NewCmd(params),
		flavor.NewCmd(params),
	)
}
