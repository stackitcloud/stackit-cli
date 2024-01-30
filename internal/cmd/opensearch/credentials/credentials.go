package credentials

import (
	"stackit/internal/cmd/opensearch/credentials/create"
	"stackit/internal/cmd/opensearch/credentials/delete"
	"stackit/internal/cmd/opensearch/credentials/describe"
	"stackit/internal/cmd/opensearch/credentials/list"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for OpenSearch credentials",
		Long:  "Provides functionality for OpenSearch credentials",
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
