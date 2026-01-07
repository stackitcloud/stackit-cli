package snapshot

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/snapshot/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/snapshot/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/snapshot/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/snapshot/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Provides functionality for SFS snapshots",
		Long:  "Provides functionality for SFS snapshots.",
		Args:  args.NoArgs,
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
}
