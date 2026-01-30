package access_token

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/delete_all"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/delete_all_expired"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs/access_token/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-token",
		Short: "Provides functionality for Logs access-tokens",
		Long:  "Provides functionality for Logs access-tokens.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(delete_all.NewCmd(params))
	cmd.AddCommand(delete_all_expired.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
