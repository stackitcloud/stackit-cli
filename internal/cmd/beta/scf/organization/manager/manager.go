package manager

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/manager/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/manager/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/manager/describe"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit beta scf organization manager' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manager",
		Short: "Manages the organization manager of STACKIT Cloud Foundry organizations",
		Long:  "Manages the organization manager of STACKIT Cloud Foundry organizations.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
}
