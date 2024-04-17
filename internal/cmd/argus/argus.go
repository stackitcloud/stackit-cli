package argus

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/plans"
	scrapeconfigs "github.com/stackitcloud/stackit-cli/internal/cmd/argus/scrape-configs"
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
	cmd.AddCommand(plans.NewCmd(p))
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(scrapeconfigs.NewCmd(p))
}
