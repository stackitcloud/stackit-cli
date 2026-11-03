package space

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/role"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/space/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit beta scf space' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Provides functionality for STACKIT Cloud Foundry spaces",
		Long:  "Provides functionality for STACKIT Cloud Foundry spaces.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(role.NewCmd(params))
}
