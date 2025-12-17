package credentials

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme/credentials/list"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for LogMe credentials",
		Long:  "Provides functionality for LogMe credentials.",
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
