package instance

import (
	"stackit/internal/cmd/opensearch/instance/create"
	"stackit/internal/cmd/opensearch/instance/delete"
	"stackit/internal/cmd/opensearch/instance/describe"
	"stackit/internal/cmd/opensearch/instance/list"
	"stackit/internal/cmd/opensearch/instance/update"
	"stackit/internal/pkg/args"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for OpenSearch instances",
		Long:  "Provides functionality for OpenSearch instances",
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
	cmd.AddCommand(update.NewCmd())
}
