package dns

import (
	recordset "github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Provides functionality for DNS",
		Long:  "Provides functionality for DNS.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(zone.NewCmd(params))
	cmd.AddCommand(recordset.NewCmd(params))
}
