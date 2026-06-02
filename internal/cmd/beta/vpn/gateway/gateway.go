package gateway

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/gateway/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "Provides functionality for VPN gateway",
		Long:  "Provides functionality for VPN gateway.",
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
	cmd.AddCommand(update.NewCmd(params))
}
