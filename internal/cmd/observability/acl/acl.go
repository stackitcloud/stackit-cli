package acl

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/acl/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/acl/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Provides functionality for ACLs of Observability instances",
		Long:  "Provides functionality for access control lists (ACLs) of Observability instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
