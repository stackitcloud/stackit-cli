package opensearch

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opensearch",
		Short: "Provides functionality for OpenSearch",
		Long:  "Provides functionality for OpenSearch.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(instance.NewCmd())
	cmd.AddCommand(plans.NewCmd())
	cmd.AddCommand(credentials.NewCmd())
}
