package config

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/config/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/config/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/config/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb-waf/config/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Provides functionality for WAF configurations of the ALB WAF",
		Long:  "Provides functionality for Web Application Firewall (WAF) configurations for application loadbalancers.",
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
