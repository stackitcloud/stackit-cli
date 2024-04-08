package zone

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "Provides functionality for DNS zones",
		Long:  "Provides functionality for DNS zones.",
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
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
}
