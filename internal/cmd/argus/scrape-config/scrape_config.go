package scrapeconfig

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config/delete"
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-config/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scrape-config",
		Short: "Provides functionality for scrape configurations in Argus",
		Long:  "Provides functionality for scrape configurations in Argus.",
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
}
