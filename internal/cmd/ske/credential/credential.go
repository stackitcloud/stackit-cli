package credential

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credential/describe"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "Provides functionality for SKE credentials",
		Long:    "Provides functionality for SKE credentials",
		Example: describe.NewCmd().Example,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(describe.NewCmd())
}
