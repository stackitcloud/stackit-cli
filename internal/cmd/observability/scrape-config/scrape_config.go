package scrapeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/describe"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/scrape-config/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scrape-config",
		Short: "Provides functionality for scrape configurations in Observability",
		Long:  "Provides functionality for scrape configurations in Observability.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(generatepayload.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
}
