package managedruleset

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set/update"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "managed-rule-set",
		Short: "Provides functionality for managed rule sets of the ALB WAF",
		Long:  "Provides functionality for managed rule sets (MRS) of the Web Application Firewall (WAF) for application loadbalancers.",
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
