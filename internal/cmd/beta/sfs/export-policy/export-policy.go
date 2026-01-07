package exportpolicy

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-policy",
		Short: "Provides functionality for SFS export policies",
		Long:  "Provides functionality for SFS export policies.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
