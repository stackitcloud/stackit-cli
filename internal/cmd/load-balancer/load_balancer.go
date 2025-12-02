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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "load-balancer",
		Aliases: []string{"lb"},
		Short:   "Provides functionality for Load Balancer",
		Long:    "Provides functionality for Load Balancer.",
		Args:    args.NoArgs,
		Run:     utils.CmdHelp,
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
	cmd.AddCommand(quota.NewCmd(params))
	cmd.AddCommand(observabilitycredentials.NewCmd(params))
	cmd.AddCommand(targetpool.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
