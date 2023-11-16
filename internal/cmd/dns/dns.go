package dns

import (
	recordset "stackit/internal/cmd/dns/record-set"
	"stackit/internal/cmd/dns/zone"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Provides functionality for DNS",
		Long:  "Provides functionality for DNS",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(zone.NewCmd())
	cmd.AddCommand(recordset.NewCmd())
}
