package zone

import (
	"stackit/internal/cmd/dns/zone/create"
	"stackit/internal/cmd/dns/zone/delete"
	"stackit/internal/cmd/dns/zone/describe"
	"stackit/internal/cmd/dns/zone/list"
	"stackit/internal/cmd/dns/zone/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "Provides functionality for DNS zones",
		Long:  "Provides functionality for DNS zones",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(update.NewCmd())
	cmd.AddCommand(delete.NewCmd())
}
