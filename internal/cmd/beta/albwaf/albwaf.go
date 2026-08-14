package albwaf

import (
	customrulegroup "github.com/stackitcloud/stackit-cli/internal/cmd/beta/albwaf/custom-rule-group"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alb-waf",
		Short: "Provides functionality for Application Load Balancer Web Application Firewall",
		Long:  "Provides functionality for Application Load Balancer Web Application Firewall.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(customrulegroup.NewCmd(params))
}
