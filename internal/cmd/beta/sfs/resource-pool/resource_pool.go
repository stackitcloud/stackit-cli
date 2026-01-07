package resourcepool

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-pool",
		Short: "Provides functionality for SFS resource pools",
		Long:  "Provides functionality for SFS resource pools.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
