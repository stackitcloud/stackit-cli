package credentials

import (
	"github.com/spf13/cobra"

	add "github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/add"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb/observability-credentials/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observability-credentials",
		Short: "Provides functionality for application loadbalancer credentials",
		Long:  "Provides functionality for application loadbalancer credentials",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(add.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
