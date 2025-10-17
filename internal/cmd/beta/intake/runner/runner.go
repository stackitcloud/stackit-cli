package runner

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/runner/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "Provides functionality for Intake Runners",
		Long:  "Provides functionality for Intake Runners.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	// Pass the params down to each action command
	cmd.AddCommand(create.NewCreateCmd(params))
	cmd.AddCommand(delete.NewDeleteCmd(params))
	cmd.AddCommand(describe.NewDescribeCmd(params))
	cmd.AddCommand(list.NewListCmd(params))
	cmd.AddCommand(update.NewUpdateCmd(params))

	return cmd
}
