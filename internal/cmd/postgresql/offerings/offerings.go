package offerings

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offerings/list"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "offerings",
	Short:   "Provides information regarding the PostgreSQL service offerings",
	Long:    "Provides information regarding the PostgreSQL service offerings",
	Example: list.Cmd.Example,
}

func init() {
	// Add all direct child commands
	Cmd.AddCommand(list.Cmd)
}
