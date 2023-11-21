package config

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/inspect"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/unset"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "CLI configuration options",
		Long:    "CLI configuration options",
		Example: fmt.Sprintf("%s\n%s\n%s", set.NewCmd().Example, inspect.NewCmd().Example, unset.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(inspect.NewCmd())
	cmd.AddCommand(set.NewCmd())
	cmd.AddCommand(unset.NewCmd())
}
