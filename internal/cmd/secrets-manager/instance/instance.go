package instance

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for Secrets Manager instances",
		Long:  "Provides functionality for Secrets Manager instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
}
