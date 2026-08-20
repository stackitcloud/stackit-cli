package albwaf

import (
	customrulegroup "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/custom-rule-group"
	managedruleset "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/managed-rule-set"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alb-waf",
		Short: "Manages the Web Application Firewall (WAF) for application loadbalancers",
		Long:  "Manage the lifecycle of Web Application Firewall (WAF) configurations for application loadbalancers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(
		managedruleset.NewCmd(params),
		customrulegroup.NewCmd(params),
	)
}
