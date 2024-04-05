package opensearch

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opensearch",
		Short: "Provides functionality for OpenSearch",
		Long:  "Provides functionality for OpenSearch.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(instance.NewCmd(p))
	cmd.AddCommand(plans.NewCmd(p))
	cmd.AddCommand(credentials.NewCmd(p))
}
