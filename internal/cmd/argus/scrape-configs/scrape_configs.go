package scrapeconfigs

import (
	generatepayload "github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-configs/generate-payload"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scrape-configs",
		Short: "Provides functionality for scrape configs in Argus.",
		Long:  "Provides functionality for scrape configurations in Argus.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(generatepayload.NewCmd(p))
}
