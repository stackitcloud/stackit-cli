package dns

import (
	recordset "github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Provides functionality for DNS",
		Long:  "Provides functionality for DNS.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(zone.NewCmd(p))
	cmd.AddCommand(recordset.NewCmd())
}
