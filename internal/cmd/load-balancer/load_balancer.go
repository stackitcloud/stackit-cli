package loadbalancer

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/quota"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"lb"},
		Short:   "Provides functionality for Load Balancer",
		Long:    "Provides functionality for Load Balancer.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(quota.NewCmd(p))
}
