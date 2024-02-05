package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for OpenSearch credentials",
		Long:  "Provides functionality for OpenSearch credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
}
