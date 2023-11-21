package postgresql

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offerings"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "postgresql",
		Short:   "Provides functionality for PostgreSQL",
		Long:    "Provides functionality for PostgreSQL",
		Example: fmt.Sprintf("%s\n%s", instance.NewCmd().Example, credential.NewCmd().Example),
	}
	addChilds(cmd)
	return cmd
}

func addChilds(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(offerings.NewCmd())
	cmd.AddCommand(credential.NewCmd())
}
