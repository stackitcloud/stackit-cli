package intake

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user"
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
	cmd.AddCommand(user.NewCmd(params))

	// Intake instance subcommands
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
}
