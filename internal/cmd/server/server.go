package server

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/backup"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/command"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/console"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/deallocate"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/log"
	machinetype "github.com/stackitcloud/stackit-cli/internal/cmd/server/machine-type"
	networkinterface "github.com/stackitcloud/stackit-cli/internal/cmd/server/network-interface"
	osUpdate "github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update"
	publicip "github.com/stackitcloud/stackit-cli/internal/cmd/server/public-ip"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/reboot"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/rescue"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/resize"
	serviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/start"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/stop"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/unrescue"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/volume"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Provides functionality for servers",
		Long:  "Provides functionality for servers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(backup.NewCmd(params))
	cmd.AddCommand(command.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(publicip.NewCmd(params))
	cmd.AddCommand(serviceaccount.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(volume.NewCmd(params))
	cmd.AddCommand(networkinterface.NewCmd(params))
	cmd.AddCommand(console.NewCmd(params))
	cmd.AddCommand(log.NewCmd(params))
	cmd.AddCommand(start.NewCmd(params))
	cmd.AddCommand(stop.NewCmd(params))
	cmd.AddCommand(reboot.NewCmd(params))
	cmd.AddCommand(deallocate.NewCmd(params))
	cmd.AddCommand(resize.NewCmd(params))
	cmd.AddCommand(rescue.NewCmd(params))
	cmd.AddCommand(unrescue.NewCmd(params))
	cmd.AddCommand(osUpdate.NewCmd(params))
	cmd.AddCommand(machinetype.NewCmd(params))
}
