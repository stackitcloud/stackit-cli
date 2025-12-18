package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for OpenSearch instances",
		Long:  "Provides functionality for OpenSearch instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
