package scrapeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/update"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scrape-config",
		Short: "Provides functionality for scrape configurations in Observability",
		Long:  "Provides functionality for scrape configurations in Observability.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(generatepayload.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
