package credential

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/list"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credential",
		Short:   "Provides functionality for PostgreSQL credentials",
		Long:    "Provides functionality for PostgreSQL credentials",
		Example: fmt.Sprintf("%s\n%s", create.NewCmd().Example, describe.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
}
