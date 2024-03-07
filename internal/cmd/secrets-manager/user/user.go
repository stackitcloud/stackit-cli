package user

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/update"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for Secrets Manager users",
		Long:  "Provides functionality for Secrets Manager users.",
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
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
