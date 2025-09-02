package intake

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/intake/runner"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit intake' command
func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intake",
		Short: "Provides functionality for STACKIT Intake",
		Long:  "Provides functionality for STACKIT Intake.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				``,
				"$ stackit intake"),
		),
		Run: utils.CmdHelp,
	}

	// Sub-commands
	cmd.AddCommand(runner.NewCmd(params))

	return cmd
}
