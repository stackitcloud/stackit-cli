package user

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake/user/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Provides functionality for Intake Users",
		Long:  "Provides functionality for Intake Users.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	// Pass the params down to each action command
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))

	return cmd
}
