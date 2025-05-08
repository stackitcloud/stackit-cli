package user

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user/update"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for Secrets Manager users",
		Long:  "Provides functionality for Secrets Manager users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
