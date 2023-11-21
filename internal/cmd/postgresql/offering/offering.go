package offering

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresql/offering/list"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "offering",
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
