package loadbalancer

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/list"
	observabilitycredentials "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/observability-credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/quota"
	targetpool "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/update"

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
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(generatepayload.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(quota.NewCmd(p))
	cmd.AddCommand(observabilitycredentials.NewCmd(p))
	cmd.AddCommand(targetpool.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
