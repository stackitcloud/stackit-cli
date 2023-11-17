package credential

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential/list"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "credential",
	Short:   "Provides functionality for PostgreSQL credentials",
	Long:    "Provides functionality for PostgreSQL credentials",
	Example: fmt.Sprintf("%s\n%s", create.Cmd.Example, describe.Cmd.Example),
}

func init() {
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(delete.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(list.Cmd)
}
