package offerings

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offerings/list"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "offerings",
		Short:   "Provides information regarding the PostgreSQL service offerings",
		Long:    "Provides information regarding the PostgreSQL service offerings",
		Example: list.NewCmd().Example,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
}
