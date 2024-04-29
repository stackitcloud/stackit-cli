package argus

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/plans"
	scrapeconfig "github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argus",
		Short: "Provides functionality for Argus",
		Long:  "Provides functionality for Argus.",
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
