package runner

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Provides functionality for Intake Runners",
		Long:  "Provides functionality for Intake Runners.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}

	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	// Pass the params down to each action command
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
