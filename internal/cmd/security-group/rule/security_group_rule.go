package rule

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Provides functionality for security group rules",
		Long:  "Provides functionality for security group rules.",
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
