package instance

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance/update"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "instance",
		Short:   "Provides functionality for PostgreSQL instance",
		Long:    "Provides functionality for PostgreSQL instance",
		Example: fmt.Sprintf("%s\n%s", create.NewCmd().Example, list.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
