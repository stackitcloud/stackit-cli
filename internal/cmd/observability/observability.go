package observability

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/plans"
	scrapeconfig "github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observability",
		Short: "Provides functionality for Observability",
		Long:  "Provides functionality for Observability.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(grafana.NewCmd(p))
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
	cmd.AddCommand(scrapeconfig.NewCmd(p))
	cmd.AddCommand(plans.NewCmd(p))
}
