package organization

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit beta scf organization' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "organization",
		Short: "Provides functionality for STACKIT Cloud Foundry organizations",
		Long:  "Provides functionality for STACKIT Cloud Foundry organizations.",
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
}
