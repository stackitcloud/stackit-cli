package projectlock

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/project-lock/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/project-lock/lock"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/project-lock/unlock"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project-lock",
		Short: "Provides functionality for SFS project locks",
		Long:  "Provides functionality for SFS project locks.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(lock.NewCmd(params))
	cmd.AddCommand(unlock.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
