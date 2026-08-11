package credentials

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/valkey/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for Valkey credentials",
		Long:  "Provides functionality for Valkey credentials.",
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
}
