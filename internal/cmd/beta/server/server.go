package server

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/command"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/console"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/deallocate"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/log"
	machinetype "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/machine-type"
	networkinterface "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/network-interface"
	osUpdate "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update"
	publicip "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/public-ip"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/reboot"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/rescue"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/resize"
	serviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/start"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/stop"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/unrescue"
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
	cmd.AddCommand(serviceaccount.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(volume.NewCmd(p))
	cmd.AddCommand(networkinterface.NewCmd(p))
	cmd.AddCommand(console.NewCmd(p))
	cmd.AddCommand(log.NewCmd(p))
	cmd.AddCommand(start.NewCmd(p))
	cmd.AddCommand(stop.NewCmd(p))
	cmd.AddCommand(reboot.NewCmd(p))
	cmd.AddCommand(deallocate.NewCmd(p))
	cmd.AddCommand(resize.NewCmd(p))
	cmd.AddCommand(rescue.NewCmd(p))
	cmd.AddCommand(unrescue.NewCmd(p))
	cmd.AddCommand(osUpdate.NewCmd(p))
	cmd.AddCommand(machinetype.NewCmd(p))
}
