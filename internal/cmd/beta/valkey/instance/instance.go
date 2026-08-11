package instance

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/instance/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Provides functionality for Valkey instances",
		Long:  "Provides functionality for Valkey instances.",
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
