package recordset

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-set",
		Short: "Provides functionality for DNS record set",
		Long:  "Provides functionality for DNS record set.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
