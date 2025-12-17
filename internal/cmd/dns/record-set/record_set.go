package recordset

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-set",
		Short: "Provides functionality for DNS record set",
		Long:  "Provides functionality for DNS record set.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
