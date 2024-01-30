package opensearch

import (
	"stackit/internal/cmd/opensearch/credentials"
	"stackit/internal/cmd/opensearch/instance"
	"stackit/internal/cmd/opensearch/plans"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opensearch",
		Short: "Provides functionality for OpenSearch",
		Long:  "Provides functionality for OpenSearch",
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
