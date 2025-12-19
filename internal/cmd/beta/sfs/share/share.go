package share

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs/share/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Provides functionality for SFS shares",
		Long:  "Provides functionality for SFS shares.",
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
