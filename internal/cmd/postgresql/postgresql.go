package postgresql

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offering"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresql",
		Short:   "Provides functionality for PostgreSQL",
		Long:    "Provides functionality for PostgreSQL",
		Example: fmt.Sprintf("%s\n%s", instance.NewCmd().Example, credential.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(offering.NewCmd())
	cmd.AddCommand(credential.NewCmd())
}
