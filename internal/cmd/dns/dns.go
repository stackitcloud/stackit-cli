package dns

import (
	"fmt"

	recordset "github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dns",
		Short:   "Provides functionality for DNS",
		Long:    "Provides functionality for DNS",
		Example: fmt.Sprintf("%s\n%s", zone.NewCmd().Example, recordset.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(zone.NewCmd())
	cmd.AddCommand(recordset.NewCmd())
}
