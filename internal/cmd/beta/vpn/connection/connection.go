package connection

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/connection/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/connection/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/connection/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/connection/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection",
		Short: "Provides functionality for VPN connections",
		Long:  "Provides functionality for VPN connections.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
