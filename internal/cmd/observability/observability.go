package observability

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/plans"
	scrapeconfig "github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "observability",
		Short: "Provides functionality for Observability",
		Long:  "Provides functionality for Observability.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(grafana.NewCmd(params))
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
	cmd.AddCommand(scrapeconfig.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
}
