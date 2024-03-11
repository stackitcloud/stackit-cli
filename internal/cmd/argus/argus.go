package argus

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argus",
		Short: "Provides functionality for Argus",
		Long:  "Provides functionality for Argus.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(plans.NewCmd())
}
