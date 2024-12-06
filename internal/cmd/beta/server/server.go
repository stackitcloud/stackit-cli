package server

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/list"
	networkinterface "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/network-interface"
	publicip "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/public-ip"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/volume"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Provides functionality for servers",
		Long:  "Provides functionality for servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(backup.NewCmd(p))
	cmd.AddCommand(command.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(publicip.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(volume.NewCmd(p))
	cmd.AddCommand(networkinterface.NewCmd(p))
}
