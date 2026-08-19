package customrulegroup

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-rule-group",
		Short: "Provides functionality for custom rule groups of the ALB WAF",
		Long:  "Provides functionality for custom rule groups (CRG) of the Web Application Firewall (WAF) for application loadbalancers.",
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
	cmd.AddCommand(generatepayload.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
