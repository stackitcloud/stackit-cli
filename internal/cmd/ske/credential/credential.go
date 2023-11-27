package credential

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credential/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credential/rotate"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "Provides functionality for SKE credentials",
		Long:    "Provides functionality for SKE credentials",
		Example: fmt.Sprintf("%s\n%s", describe.NewCmd().Example, rotate.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(rotate.NewCmd())
}
