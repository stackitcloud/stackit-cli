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

var Cmd = &cobra.Command{
	Use:     "instance",
	Short:   "Provides functionality for PostgreSQL instance",
	Long:    "Provides functionality for PostgreSQL instance",
	Example: fmt.Sprintf("%s\n%s", create.Cmd.Example, list.Cmd.Example),
}

func init() {
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(delete.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(update.Cmd)
}
