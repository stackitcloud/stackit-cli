package key

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Provides functionality regarding service account keys",
		Long:  "Provides functionality regarding service account keys",
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
