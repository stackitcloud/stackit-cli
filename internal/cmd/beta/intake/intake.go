package intake

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit intake' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intake",
		Short: "Provides functionality for intake",
		Long:  "Provides functionality for intake.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(runner.NewCmd(params))

	// Intake instance subcommands
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
}
