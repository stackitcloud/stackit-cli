package osupdate

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/enable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/server/os-update/schedule"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "os-update",
		Short: "Provides functionality for managed server updates",
		Long:  "Provides functionality for managed server updates.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(disable.NewCmd(p))
	cmd.AddCommand(schedule.NewCmd(p))
}
