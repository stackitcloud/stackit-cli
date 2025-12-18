package opensearch

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/plans"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opensearch",
		Short: "Provides functionality for OpenSearch",
		Long:  "Provides functionality for OpenSearch.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(credentials.NewCmd(params))
}
