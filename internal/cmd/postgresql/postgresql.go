package postgresql

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/credential"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offerings"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "postgresql",
	Short:   "Provides functionality for PostgreSQL",
	Long:    "Provides functionality for PostgreSQL",
	Example: fmt.Sprintf("%s\n%s", instance.Cmd.Example, credential.Cmd.Example),
}

func init() {
	Cmd.AddCommand(instance.Cmd)
	Cmd.AddCommand(offerings.Cmd)
	Cmd.AddCommand(credential.Cmd)
}
