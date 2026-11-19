package quota

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/quota/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/quota/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/scf/organization/quota/set"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// NewCmd creates the 'stackit beta scf organization quota' command
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Manages quotas of STACKIT Cloud Foundry organizations",
		Long:  "Manages quotas of STACKIT Cloud Foundry organizations.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(set.NewCmd(params))
}
