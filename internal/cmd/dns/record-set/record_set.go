package recordset

import (
	"stackit/internal/cmd/dns/record-set/create"
	"stackit/internal/cmd/dns/record-set/delete"
	"stackit/internal/cmd/dns/record-set/describe"
	"stackit/internal/cmd/dns/record-set/list"
	"stackit/internal/cmd/dns/record-set/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-set",
		Short: "Provides functionality for DNS record set",
		Long:  "Provides functionality for DNS record set",
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
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
