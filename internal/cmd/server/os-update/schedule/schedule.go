package schedule

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/create"
	del "github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/os-update/schedule/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Provides functionality for Server os-update Schedule",
		Long:  "Provides functionality for Server os-update Schedule.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(del.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
