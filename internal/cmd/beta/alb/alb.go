package alb

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/template"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alb",
		Short: "Manages application loadbalancers",
		Long:  "Manage the lifecycle of application loadbalancers.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(
		list.NewCmd(p),
		template.NewCmd(p),
		create.NewCmd(p),
		update.NewCmd(p),
		describe.NewCmd(p),
		delete.NewCmd(p),
	)
}
