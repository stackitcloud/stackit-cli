package vpn

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/connection"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/plans"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/quotas"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpn",
		Short: "Provides functionality for VPN",
		Long:  "Provides functionality for VPN.",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(connection.NewCmd(params))
	cmd.AddCommand(gateway.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(quotas.NewCmd(params))
}
