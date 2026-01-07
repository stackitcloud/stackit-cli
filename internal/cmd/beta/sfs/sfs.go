package sfs

import (
	exportpolicy "github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/export-policy"
	performanceclass "github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/performance-class"
	resourcepool "github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/resource-pool"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/snapshot"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sfs",
		Short: "Provides functionality for SFS (stackit file storage)",
		Long:  "Provides functionality for SFS (stackit file storage).",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(resourcepool.NewCmd(params))
	cmd.AddCommand(share.NewCmd(params))
	cmd.AddCommand(exportpolicy.NewCmd(params))
	cmd.AddCommand(snapshot.NewCmd(params))
	cmd.AddCommand(performanceclass.NewCmd(params))
}
