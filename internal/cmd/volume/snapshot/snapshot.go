package snapshot

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/snapshot/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/snapshot/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/snapshot/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/snapshot/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume/snapshot/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Provides functionality for snapshots",
		Long:  "Provides functionality for snapshots.",
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
	cmd.AddCommand(update.NewCmd(params))
}
